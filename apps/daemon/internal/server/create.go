package server

import (
	"encoding/json"
	"fmt"
	"log"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
	"panelium/daemon/internal/sync"
	"panelium/proto_gen_go/daemon"
	"slices"
)

func CreateServer(sid string, ownerId string, userIds []string, allocations []model.ServerAllocation, resourceLimit model.ResourceLimit, dockerImage string, bid string) (*model.Server, error) {
	err := sync.SyncBlueprints()
	if err != nil {
		log.Printf("failed to sync blueprints: %v", err)
	}

	blueprint := model.Blueprint{}
	tx := db.Instance().First(&blueprint, "bid = ?", bid)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, fmt.Errorf("failed to find blueprint with ID %s: %w", bid, tx.Error)
	}

	var dockerImages []string
	err = json.Unmarshal(blueprint.DockerImages, &dockerImages)
	if err != nil {
		return nil, fmt.Errorf("failed to scan docker images from blueprint: %w", err)
	}

	if !slices.Contains(dockerImages, dockerImage) {
		return nil, fmt.Errorf("docker image %s is not allowed by the blueprint %s", dockerImage, bid)
	}

	server := model.Server{
		SID:           sid,
		OwnerID:       ownerId,
		Status:        daemon.ServerStatusType_SERVER_STATUS_TYPE_INSTALLING,
		ResourceLimit: resourceLimit,
		DockerImage:   dockerImage,
		BID:           bid,
	}
	tx = db.Instance().Create(&server)
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Printf("failed to create server: %v\n", tx.Error)
		return nil, fmt.Errorf("failed to create server: %w", tx.Error)
	}

	if len(userIds) > 0 {
		for _, userId := range userIds {
			serverUser := model.ServerUser{
				SID: server.SID,
				UID: userId,
			}
			if err := db.Instance().Create(&serverUser).Error; err != nil {
				return nil, fmt.Errorf("failed to create server user: %w", err)
			}
		}
	}

	for _, allocation := range allocations {
		if allocation.Port > 65535 || allocation.Port < 1024 {
			return nil, fmt.Errorf("port %d is out of range (1024-65535)", allocation.Port)
		}

		serverAllocation := model.ServerAllocation{
			IP:   allocation.IP,
			Port: allocation.Port,
			SID:  server.SID,
		}
		if err := db.Instance().Create(&serverAllocation).Error; err != nil {
			return nil, err
		}
	}

	go func() {
		err := Install(server.SID)
		if err != nil {
			log.Printf("failed to install server %s: %v\n", server.SID, err)
			return
		}

		err = Start(server.SID) // TODO: maybe move to install?
		if err != nil {
			log.Printf("failed to start server %s: %v\n", server.SID, err)
			return
		}
	}()

	return &server, nil
}
