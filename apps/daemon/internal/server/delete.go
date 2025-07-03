package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"panelium/daemon/internal/docker"
)

func DeleteServer(sid string) error {
	dbErr := yeetDbServer(sid)
	if dbErr != nil {
		fmt.Printf("failed to delete server %s from database: %v\n", sid, dbErr)
	}
	crErr := docker.Instance().ContainerRemove(context.Background(), sid, container.RemoveOptions{
		Force: true,
	})
	if crErr != nil {
		fmt.Printf("failed to remove server container %s: %v\n", sid, crErr)
	}
	volErr := docker.Instance().VolumeRemove(context.Background(), sid, true)
	if volErr != nil {
		fmt.Printf("failed to remove server volume %s: %v\n", sid, volErr)
	}
	netErr := docker.Instance().NetworkRemove(context.Background(), sid)
	if netErr != nil {
		fmt.Printf("failed to remove server network %s: %v\n", sid, netErr)
	}

	if dbErr != nil || crErr != nil || volErr != nil || netErr != nil {
		return errors.Join(errors.New("failed to delete server completely"), errors.Join(dbErr, crErr, volErr, netErr))
	}

	return nil
}
