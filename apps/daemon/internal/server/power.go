package server

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"panelium/common/util"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/model"
	"panelium/proto_gen_go/daemon"
	"time"
)

func Start(s *model.Server) error {
	if !s.ContainerExists {
		return fmt.Errorf("server %s does not have a container", s.SID)
	}

	ci, err := docker.Instance().ContainerInspect(context.Background(), s.SID)
	if err != nil {
		return fmt.Errorf("failed to inspect server container %s: %v", s.SID, err)
	}
	if ci.State.Running {
		return nil
	}

	err = docker.Instance().ContainerStart(context.Background(), s.SID, container.StartOptions{})
	if err != nil {
		fmt.Printf("failed to start server container %s: %v\n", s.SID, err)
	}
	s.Status = daemon.ServerStatusType_SERVER_STATUS_TYPE_STARTING
	s.TimestampStart = time.Now()
	if err := db.Instance().Save(s).Error; err != nil {
		return err
	}

	// TODO: start a goroutine to wait for the server to be fully started and update the status

	return nil
}

var stopTimeout = -1
var killTimeout = 0

func Stop(s *model.Server, kill bool) error {
	if !s.ContainerExists {
		return fmt.Errorf("server %s does not have a container", s.SID)
	}

	// might split into Stop and Kill functions later with kill using ContainerKill
	err := docker.Instance().ContainerStop(context.Background(), s.SID, container.StopOptions{
		Timeout: util.IfElse(kill, &stopTimeout, &killTimeout),
	})
	if err != nil {
		fmt.Printf("failed to stop server container %s: %v\n", s.SID, err)
		return err
	}
	s.Status = daemon.ServerStatusType_SERVER_STATUS_TYPE_STOPPING
	if err := db.Instance().Save(s).Error; err != nil {
		return err
	}

	go func() {
		statusCh, errCh := docker.Instance().ContainerWait(context.Background(), s.SID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				fmt.Printf("error waiting for server container %s to stop: %v\n", s.SID, err)
			}
		case status := <-statusCh:
			if status.StatusCode != 0 {
				fmt.Printf("server container %s stopped with non-zero status code: %d\n", s.SID, status.StatusCode)
			}

			s.Status = daemon.ServerStatusType_SERVER_STATUS_TYPE_OFFLINE
			s.OfflineReason = util.IfElse(kill, daemon.ServerOfflineReason_SERVER_OFFLINE_REASON_KILLED, daemon.ServerOfflineReason_SERVER_OFFLINE_REASON_STOPPED)
			s.TimestampEnd = time.Now()
			if err := db.Instance().Save(s).Error; err != nil {
				fmt.Printf("failed to update server status after stop: %v\n", err)
			}
		}
	}()

	return nil
}

func Restart(s *model.Server) error {
	if !s.ContainerExists {
		return fmt.Errorf("server %s does not have a container", s.SID)
	}

	err := docker.Instance().ContainerRestart(context.Background(), s.SID, container.StopOptions{
		Timeout: &stopTimeout,
	})
	if err != nil {
		fmt.Printf("failed to restart server container %s: %v\n", s.SID, err)
		return err
	}
	s.Status = daemon.ServerStatusType_SERVER_STATUS_TYPE_STOPPING
	s.TimestampEnd = time.Now()
	if err := db.Instance().Save(s).Error; err != nil {
		return err
	}

	go func() {
		statusCh, errCh := docker.Instance().ContainerWait(context.Background(), s.SID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				fmt.Printf("error waiting for server container %s to stop: %v\n", s.SID, err)

				s.OfflineReason = daemon.ServerOfflineReason_SERVER_OFFLINE_REASON_ERROR
				if err := db.Instance().Save(s).Error; err != nil {
					fmt.Printf("failed to update server status after restart error: %v\n", err)
				}
			}
		case status := <-statusCh:
			if status.StatusCode != 0 {
				fmt.Printf("server container %s stopped with non-zero status code: %d\n", s.SID, status.StatusCode)

				s.OfflineReason = daemon.ServerOfflineReason_SERVER_OFFLINE_REASON_ERROR
				if err := db.Instance().Save(s).Error; err != nil {
					fmt.Printf("failed to update server status after restart error: %v\n", err)
				}
			}

			s.Status = daemon.ServerStatusType_SERVER_STATUS_TYPE_STARTING
			s.TimestampStart = time.Now()
			if err := db.Instance().Save(s).Error; err != nil {
				fmt.Printf("failed to update server status after stop: %v\n", err)
			}
		}
	}()

	// TODO: start a goroutine to wait for the server to be fully restarted and update the status

	return nil
}
