package server

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *ServerServiceHandler) CreateServer(
	ctx context.Context,
	req *connect.Request[proto_gen_go.CreateServerRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	server := model.Server{
		SID:    req.Msg.ServerId,
		Status: proto_gen_go.ServerStatusType_SERVER_STATUS_TYPE_INSTALLING,
		ResourceLimit: model.ResourceLimit{
			CPU:     req.Msg.ResourceLimit.Cpu,
			RAM:     req.Msg.ResourceLimit.Ram,
			SWAP:    req.Msg.ResourceLimit.Swap,
			Storage: req.Msg.ResourceLimit.Storage,
		},
		DockerImage: req.Msg.DockerImage,
		BID:         req.Msg.BlueprintId,
	}
	if err := db.Instance().Create(&server).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, allocation := range req.Msg.Allocations {
		if allocation.Port > 65535 || allocation.Port < 1024 {
			if err := db.Instance().Model(&model.ServerAllocation{}).Where("sid = ?", server.ID).Delete(&model.ServerAllocation{}).Error; err != nil {
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to rollback server allocation creation: %w", err))
			}

			if err := db.Instance().Delete(&server).Error; err != nil {
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to rollback server creation: %w", err))
			}

			return nil, connect.NewError(connect.CodeInvalidArgument, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("port %d is out of range (1024-65535)", allocation.Port)))
		}

		serverAllocation := model.ServerAllocation{
			IP:       allocation.Ip,
			Port:     uint16(allocation.Port),
			ServerID: server.ID,
		}
		if err := db.Instance().Create(&serverAllocation).Error; err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// TODO: docker

	return nil, nil
}
