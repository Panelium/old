package server

import (
	"fmt"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
	"panelium/proto_gen_go"
)

func CreateServer(sid string, allocations []model.ServerAllocation, resourceLimit model.ResourceLimit, dockerImage string, bid string) (*model.Server, error) {
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

	// TODO: docker

	return &server, nil
}
