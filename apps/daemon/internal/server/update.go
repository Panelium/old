package server

import (
	"fmt"
	"log"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
	"panelium/daemon/internal/sync"
	"slices"
)

func UpdateServer(sid string, userIds *[]string, allocations *[]model.ServerAllocation, resourceLimit *model.ResourceLimit, dockerImage *string, bid *string) error {
	err := sync.SyncBlueprints()
	if err != nil {
		log.Printf("failed to sync blueprints: %v", err)
	}

	if userIds != nil {
		tx := db.Instance().Delete(&model.ServerUser{}, "sid = ?", sid)
		if tx.Error != nil {
			return fmt.Errorf("failed to delete existing server users: %w", tx.Error)
		}

		for _, userId := range *userIds {
			if userId == "" {
				return fmt.Errorf("user ID cannot be empty")
			}
			serverUser := model.ServerUser{
				SID: sid,
				UID: userId,
			}
			if err := db.Instance().Create(&serverUser).Error; err != nil {
				return fmt.Errorf("failed to create server user: %w", err)
			}
		}
	}
	if allocations != nil {
		tx := db.Instance().Delete(&model.ServerAllocation{}, "sid = ?", sid)
		if tx.Error != nil {
			return fmt.Errorf("failed to delete existing server allocations: %w", tx.Error)
		}

		for _, allocation := range *allocations {
			if allocation.Port > 65535 || allocation.Port < 1024 {
				return fmt.Errorf("port %d is out of range (1024-65535)", allocation.Port)
			}

			serverAllocation := model.ServerAllocation{
				IP:   allocation.IP,
				Port: allocation.Port,
				SID:  sid,
			}
			if err := db.Instance().Create(&serverAllocation).Error; err != nil {
				return err
			}
		}
	}
	if resourceLimit != nil {
		tx := db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Updates(model.Server{
			ResourceLimit: *resourceLimit,
		})
		if tx.Error != nil {
			return fmt.Errorf("failed to update resource limit: %w", tx.Error)
		}
	}
	if bid != nil {
		blueprint := model.Blueprint{}
		tx := db.Instance().First(&blueprint, "s.BID = ?", bid)
		if tx.Error != nil || tx.RowsAffected == 0 {
			return fmt.Errorf("failed to find blueprint with ID %s: %w", bid, tx.Error)
		}

		var dockerImages []string
		err := blueprint.DockerImages.Scan(&dockerImages)
		if err != nil {
			return fmt.Errorf("failed to scan docker images from blueprint: %w", err)
		}

		if dockerImage == nil {
			server := model.Server{}
			tx = db.Instance().First(&server, "s.SID = ?", sid)
			if tx.Error != nil || tx.RowsAffected == 0 {
				return fmt.Errorf("failed to find server with ID %s: %w", sid, tx.Error)
			}
			dockerImage = &server.DockerImage
		}

		if slices.Contains(dockerImages, *dockerImage) {
			return fmt.Errorf("docker image %s is not allowed by the blueprint %s", *dockerImage, bid)
		}
	} else if dockerImage != nil {
		server := model.Server{}
		tx := db.Instance().First(&server, "s.SID = ?", sid)
		if tx.Error != nil || tx.RowsAffected == 0 {
			return fmt.Errorf("failed to find server with ID %s: %w", sid, tx.Error)
		}

		blueprint := model.Blueprint{}
		tx = db.Instance().First(&blueprint, "s.BID = ?", server.BID)
		if tx.Error != nil || tx.RowsAffected == 0 {
			return fmt.Errorf("failed to find blueprint with ID %s: %w", server.BID, tx.Error)
		}

		var dockerImages []string
		err := blueprint.DockerImages.Scan(&dockerImages)
		if err != nil {
			return fmt.Errorf("failed to scan docker images from blueprint: %w", err)
		}
		if !slices.Contains(dockerImages, *dockerImage) {
			return fmt.Errorf("docker image %s is not allowed by the blueprint %s", *dockerImage, server.BID)
		}

		tx = db.Instance().Model(&model.Server{}).Where("sid = ?", sid).Update("docker_image", *dockerImage)
		if tx.Error != nil || tx.RowsAffected == 0 {
			return fmt.Errorf("failed to update docker image: %w", tx.Error)
		}
	}

	server := model.Server{}
	tx := db.Instance().First(&server, "s.SID = ?", sid)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return fmt.Errorf("failed to find server with ID %s: %w", sid, tx.Error)
	}

	go func() {
		err := Install(&server)
		if err != nil {
			log.Printf("failed to install server %s: %v\n", server.SID, err)
			return
		}

		err = Start(&server) // TODO: maybe move to install?
		if err != nil {
			log.Printf("failed to start server %s: %v\n", server.SID, err)
			return
		}
	}()

	return nil
}
