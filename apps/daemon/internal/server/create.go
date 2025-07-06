package server

import (
	"fmt"
	"log"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
	"panelium/proto_gen_go/daemon"
	"slices"
)

func yeetDbServer(sid string) error {
	if err := db.Instance().Delete(&model.ServerUser{}, "sid = ?", sid).Error; err != nil {
		return fmt.Errorf("failed to delete server users: %w", err)
	}

	if err := db.Instance().Delete(&model.ServerAllocation{}, "sid = ?", sid).Error; err != nil {
		return fmt.Errorf("failed to delete server allocations: %w", err)
	}

	if err := db.Instance().Delete(&model.Server{}, "sid = ?", sid).Error; err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	return nil
}

func CreateServer(sid string, ownerId string, userIds []string, allocations []model.ServerAllocation, resourceLimit model.ResourceLimit, dockerImage string, bid string) (*model.Server, error) {
	blueprint := model.Blueprint{}
	tx := db.Instance().First(&blueprint, "s.BID = ?", bid)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, fmt.Errorf("failed to find blueprint with ID %s: %w", bid, tx.Error)
	}

	var dockerImages []string
	err := blueprint.DockerImages.Scan(dockerImages)
	if err != nil {
		return nil, fmt.Errorf("failed to scan docker images from blueprint: %w", err)
	}

	if slices.Contains(dockerImages, dockerImage) {
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
	if err := db.Instance().Create(&server).Error; err != nil {
		return nil, err
	}

	if len(userIds) > 0 {
		for _, userId := range userIds {
			serverUser := model.ServerUser{
				SID: server.SID,
				UID: userId,
			}
			if err := db.Instance().Create(&serverUser).Error; err != nil {
				if rollbackErr := yeetDbServer(server.SID); rollbackErr != nil {
					return nil, fmt.Errorf("failed to rollback server creation: %w", rollbackErr)
				}
				return nil, fmt.Errorf("failed to create server user: %w", err)
			}
		}
	}

	for _, allocation := range allocations {
		if allocation.Port > 65535 || allocation.Port < 1024 {
			if err := db.Instance().Model(&model.ServerAllocation{}).Where("sid = ?", server.SID).Delete(&model.ServerAllocation{}).Error; err != nil {
				return nil, fmt.Errorf("failed to rollback server allocation creation: %w", err)
			}

			if err := db.Instance().Delete(&server).Error; err != nil {
				return nil, fmt.Errorf("failed to rollback server creation: %w", err)
			}

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
		err := Install(&server)
		if err != nil {
			log.Printf("failed to install server %s: %v\n", server.SID, err)
			if rollbackErr := yeetDbServer(server.SID); rollbackErr != nil {
				log.Printf("failed to rollback server creation: %v\n", rollbackErr)
			}
			return
		}

		err = Start(&server) // TODO: maybe move to install?
		if err != nil {
			log.Printf("failed to start server %s: %v\n", server.SID, err)
			return
		}
	}()

	return &server, nil
}
