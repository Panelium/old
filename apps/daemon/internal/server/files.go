package server

import (
	"context"
	"fmt"
	"os"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/model"
)

func rootDirectory(sid string) (string, error) {
	//check if the server ID is valid
	tx := db.Instance().First(&model.Server{}, "sid = ?", sid)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return "", fmt.Errorf("server not found")
	}

	vol, err := docker.Instance().VolumeInspect(context.Background(), fmt.Sprint("server_", sid))
	if err != nil {
		return "", fmt.Errorf("failed to inspect volume for server")
	}

	return vol.Mountpoint, nil
}

func GetRoot(sid string) (*os.Root, error) {
	rootPath, err := rootDirectory(sid)
	if err != nil {
		return nil, err
	}

	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open server root directory")
	}

	return root, nil
}
