package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"log"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/model"
)

func DeleteServer(sid string, force bool) error {
	var dbErr error
	tx := db.Instance().Delete(&model.Server{}, "sid = ?", sid)
	if tx.Error != nil {
		log.Printf("failed to delete server %s from database: %v\n", sid, tx.Error)
		dbErr = tx.Error
	}

	crErr := docker.Instance().ContainerRemove(context.Background(), fmt.Sprint("server_", sid), container.RemoveOptions{
		Force: force,
	})
	if crErr != nil {
		log.Printf("failed to remove server container %s: %v\n", sid, crErr)
	}
	volErr := docker.Instance().VolumeRemove(context.Background(), fmt.Sprint("server_", sid), force)
	if volErr != nil {
		log.Printf("failed to remove server volume %s: %v\n", sid, volErr)
	}
	netErr := docker.Instance().NetworkRemove(context.Background(), fmt.Sprint("server_", sid))
	if netErr != nil {
		log.Printf("failed to remove server network %s: %v\n", sid, netErr)
	}

	if dbErr != nil || crErr != nil || volErr != nil || netErr != nil {
		return errors.Join(errors.New("failed to delete server completely"), errors.Join(dbErr, crErr, volErr, netErr))
	}

	return nil
}
