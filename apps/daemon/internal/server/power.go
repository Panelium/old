package server

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"log"
	"panelium/common/util"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/model"
	"panelium/proto_gen_go/daemon"
	"time"
)

func Start(sid string) error {
	var s model.Server
	tx := db.Instance().First(&s, "sid = ?", sid)
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Printf("err: %v\n", tx.Error)
		return fmt.Errorf("failed to find server with ID %s: %w", sid, tx.Error)
	}

	if !s.ContainerExists {
		return fmt.Errorf("server %s does not have a container", s.SID)
	}

	ci, err := docker.Instance().ContainerInspect(context.Background(), fmt.Sprint("server_", s.SID))
	if err != nil {
		return fmt.Errorf("failed to inspect server container %s: %v", s.SID, err)
	}
	if ci.State.Running {
		return nil
	}

	err = docker.Instance().ContainerStart(context.Background(), fmt.Sprint("server_", s.SID), container.StartOptions{})
	if err != nil {
		log.Printf("failed to start server container %s: %v\n", s.SID, err)
	}

	tx = db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Updates(model.Server{
		Status:         daemon.ServerStatusType_SERVER_STATUS_TYPE_ONLINE, // TODO: set to starting
		TimestampStart: time.Now(),
	})
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Printf("failed to update server status to starting: %v\n", tx.Error)
		return fmt.Errorf("failed to update server status to starting: %w", tx.Error)
	}

	// TODO: start a goroutine to wait for the server to be fully started and update the status

	return nil
}

var stopTimeout = -1
var killTimeout = 0

func Stop(sid string, kill bool) error {
	var s model.Server
	tx := db.Instance().First(&s, "sid = ?", sid)
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Printf("err: %v\n", tx.Error)
		return fmt.Errorf("failed to find server with ID %s: %w", sid, tx.Error)
	}

	if !s.ContainerExists {
		return fmt.Errorf("server %s does not have a container", s.SID)
	}

	// might split into Stop and Kill functions later with kill using ContainerKill
	err := docker.Instance().ContainerStop(context.Background(), fmt.Sprint("server_", s.SID), container.StopOptions{
		Timeout: util.IfElse(kill, &stopTimeout, &killTimeout),
	})
	if err != nil {
		log.Printf("failed to stop server container %s: %v\n", s.SID, err)
		return err
	}

	tx = db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Updates(model.Server{
		Status: daemon.ServerStatusType_SERVER_STATUS_TYPE_STOPPING,
	})
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Printf("failed to update server status to stopping: %v\n", tx.Error)
		return fmt.Errorf("failed to update server status to stopping: %w", tx.Error)
	}

	go func() {
		statusCh, errCh := docker.Instance().ContainerWait(context.Background(), fmt.Sprint("server_", s.SID), container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				log.Printf("error waiting for server container %s to stop: %v\n", s.SID, err)
			}
		case status := <-statusCh:
			if status.StatusCode != 0 {
				log.Printf("server container %s stopped with non-zero status code: %d\n", s.SID, status.StatusCode)
			}

			tx = db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Updates(model.Server{
				Status:        daemon.ServerStatusType_SERVER_STATUS_TYPE_OFFLINE,
				OfflineReason: util.IfElse(kill, daemon.ServerOfflineReason_SERVER_OFFLINE_REASON_KILLED, daemon.ServerOfflineReason_SERVER_OFFLINE_REASON_STOPPED),
				TimestampEnd:  time.Now(),
			})
			if tx.Error != nil || tx.RowsAffected == 0 {
				log.Printf("failed to update server status to offline: %v\n", tx.Error)
			}
		}
	}()

	return nil
}

func Restart(sid string) error {
	var s model.Server
	tx := db.Instance().First(&s, "sid = ?", sid)
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Printf("err: %v\n", tx.Error)
		return fmt.Errorf("failed to find server with ID %s: %w", sid, tx.Error)
	}

	if !s.ContainerExists {
		return fmt.Errorf("server %s does not have a container", s.SID)
	}

	err := docker.Instance().ContainerRestart(context.Background(), s.SID, container.StopOptions{
		Timeout: &stopTimeout,
	})
	if err != nil {
		log.Printf("failed to restart server container %s: %v\n", s.SID, err)
		return err
	}

	tx = db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Updates(model.Server{
		Status:       daemon.ServerStatusType_SERVER_STATUS_TYPE_STOPPING,
		TimestampEnd: time.Now(),
	})
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Printf("failed to update server status to stopping: %v\n", tx.Error)
		return fmt.Errorf("failed to update server status to stopping: %w", tx.Error)
	}

	go func() {
		statusCh, errCh := docker.Instance().ContainerWait(context.Background(), s.SID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				log.Printf("error waiting for server container %s to stop: %v\n", s.SID, err)

				tx = db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Updates(model.Server{
					OfflineReason: daemon.ServerOfflineReason_SERVER_OFFLINE_REASON_ERROR,
				})
				if tx.Error != nil || tx.RowsAffected == 0 {
					log.Printf("failed to update server status after restart error: %v\n", tx.Error)
				}
			}
		case status := <-statusCh:
			if status.StatusCode != 0 {
				log.Printf("server container %s stopped with non-zero status code: %d\n", s.SID, status.StatusCode)

				tx = db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Updates(model.Server{
					OfflineReason: daemon.ServerOfflineReason_SERVER_OFFLINE_REASON_ERROR,
				})
				if tx.Error != nil || tx.RowsAffected == 0 {
					log.Printf("failed to update server status after restart error: %v\n", tx.Error)
				}
			}

			tx = db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Updates(model.Server{
				Status:         daemon.ServerStatusType_SERVER_STATUS_TYPE_STARTING,
				TimestampStart: time.Now(),
			})
			if tx.Error != nil || tx.RowsAffected == 0 {
				log.Printf("failed to update server status to starting after restart: %v\n", tx.Error)
			}
		}
	}()

	// TODO: start a goroutine to wait for the server to be fully restarted and update the status

	return nil
}
