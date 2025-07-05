package backend

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"fmt"
	"panelium/daemon/internal/model"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/daemon"
)

func (s *BackendServiceHandler) UpdateServer(
	ctx context.Context,
	req *connect.Request[daemon.UpdateServerRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	var allocations *[]model.ServerAllocation = nil
	if req.Msg.Allocations != nil {
		allocs := make([]model.ServerAllocation, len(req.Msg.Allocations))
		for i, alloc := range req.Msg.Allocations {
			if alloc.Port < 1024 || alloc.Port > 65535 {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("port %d is out of range (1024-65535)", alloc.Port))
			}

			allocs[i] = model.ServerAllocation{
				IP:   alloc.Ip,
				Port: uint16(alloc.Port),
			}
		}

		allocations = &allocs
	}

	var userIds *[]string = nil
	if req.Msg.UserIds != nil {
		userIds = &req.Msg.UserIds
	}

	var resourceLimit *model.ResourceLimit = nil
	if req.Msg.ResourceLimit != nil {
		resourceLimit = &model.ResourceLimit{
			CPU:     req.Msg.ResourceLimit.Cpu,
			RAM:     req.Msg.ResourceLimit.Ram,
			SWAP:    req.Msg.ResourceLimit.Swap,
			Storage: req.Msg.ResourceLimit.Storage,
		}
	}

	err := server.UpdateServer(req.Msg.ServerId, userIds, allocations, resourceLimit, req.Msg.DockerImage, req.Msg.BlueprintId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to create server"))
	}

	res := connect.NewResponse(&proto_gen_go.SuccessMessage{
		Success: true,
	})

	return res, nil
}
