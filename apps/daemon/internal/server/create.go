package server

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/model"
	"panelium/proto_gen_go"
)

func yeetDbServer(sid string) error {
	if err := db.Instance().Model(&model.ServerAllocation{}).Where("sid = ?", sid).Delete(&model.ServerAllocation{}).Error; err != nil {
		return fmt.Errorf("failed to delete server allocations: %w", err)
	}

	if err := db.Instance().Delete(&model.Server{}, "sid = ?", sid).Error; err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	return nil
}

func CreateServer(sid string, allocations []model.ServerAllocation, resourceLimit model.ResourceLimit, dockerImage string, bid string) (*model.Server, error) {
	// image pull and container create are done in a goroutine to avoid blocking the main thread and prevent timeouts
	go func() {
		rc, err := docker.Instance().ImagePull(context.Background(), dockerImage, image.PullOptions{})
		if err != nil {
			if err := yeetDbServer(sid); err != nil {
				fmt.Printf("failed to rollback server creation: %v\n", err)
			}
			fmt.Printf("failed to pull docker image: %v\n", err)
			return
		}

		_, err = io.ReadAll(rc)
		if err != nil {
			if err := yeetDbServer(sid); err != nil {
				fmt.Printf("failed to rollback server creation: %v\n", err)
			}
			fmt.Printf("failed to read docker image pull response: %v\n", err)
			return
		}
		err = rc.Close()
		if err != nil {
			if err := yeetDbServer(sid); err != nil {
				fmt.Printf("failed to rollback server creation: %v\n", err)
			}
			fmt.Printf("failed to close docker image pull response: %v\n", err)
			return
		}

		cr, err := docker.Instance().ContainerCreate(context.Background(), &container.Config{}, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, sid)
		if err != nil {
			if err := yeetDbServer(sid); err != nil {
				fmt.Printf("failed to rollback server creation: %v\n", err)
			}
			fmt.Printf("failed to create docker container: %v\n", err)
		}

		if err := docker.Instance().ContainerStart(context.Background(), cr.ID, container.StartOptions{}); err != nil {
			if err := yeetDbServer(sid); err != nil {
				fmt.Printf("failed to rollback server creation: %v\n", err)
			}
			err := docker.Instance().ContainerRemove(context.Background(), cr.ID, container.RemoveOptions{Force: true})
			if err != nil {
				return
			}
			fmt.Printf("failed to start docker container: %v\n", err)
			return
		}
	}()

	server := model.Server{
		SID:           sid,
		Status:        proto_gen_go.ServerStatusType_SERVER_STATUS_TYPE_INSTALLING,
		ResourceLimit: resourceLimit,
		DockerImage:   dockerImage,
		BID:           bid,
	}
	if err := db.Instance().Create(&server).Error; err != nil {
		return nil, err
	}

	for _, allocation := range allocations {
		if allocation.Port > 65535 || allocation.Port < 1024 {
			if err := db.Instance().Model(&model.ServerAllocation{}).Where("sid = ?", server.ID).Delete(&model.ServerAllocation{}).Error; err != nil {
				return nil, fmt.Errorf("failed to rollback server allocation creation: %w", err)
			}

			if err := db.Instance().Delete(&server).Error; err != nil {
				return nil, fmt.Errorf("failed to rollback server creation: %w", err)
			}

			return nil, fmt.Errorf("port %d is out of range (1024-65535)", allocation.Port)
		}

		serverAllocation := model.ServerAllocation{
			IP:       allocation.IP,
			Port:     allocation.Port,
			ServerID: server.ID,
		}
		if err := db.Instance().Create(&serverAllocation).Error; err != nil {
			return nil, err
		}
	}

	return &server, nil
}
