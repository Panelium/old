package daemon

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/backend/internal/db"
	"panelium/backend/internal/middleware"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *DaemonServiceHandler) GetServer(
	ctx context.Context,
	req *connect.Request[proto_gen_go.SimpleIDMessage],
) (*connect.Response[backend.Server], error) {
	daemonInfoData := ctx.Value("panelium_daemon_info")
	daemonInfo, ok := daemonInfoData.(*middleware.DaemonInfo)
	if !ok || daemonInfo == nil || daemonInfo.NID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
	}

	var node *model.Node
	tx := db.Instance().First(&node, "nid = ?", daemonInfo.NID)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("node not found"))
	}

	var server *model.Server
	tx = db.Instance().First(&server, "sid = ? AND node_id = ?", req.Msg.Id, node.ID)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("server not found"))
	}

	var userIds []string
	for _, user := range server.Users {
		userIds = append(userIds, user.User.UID)
	}

	var allocations []*proto_gen_go.IPAllocation
	for _, allocation := range server.Allocations {
		if allocation.IP == "" || allocation.Port < 1024 || allocation.Port > 65535 {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid server allocation"))
		}

		allocations = append(allocations, &proto_gen_go.IPAllocation{
			Ip:   allocation.IP,
			Port: uint32(allocation.Port),
		})
	}

	serverProto := backend.Server{
		Sid:         server.SID,
		OwnerId:     server.Owner.UID,
		UserIds:     userIds,
		Allocations: allocations,
		ResourceLimit: &proto_gen_go.ResourceLimit{
			Cpu:     uint32(server.ResourceLimit.CPU),
			Ram:     uint32(server.ResourceLimit.RAM),
			Swap:    uint32(server.ResourceLimit.SWAP),
			Storage: uint32(server.ResourceLimit.Storage),
		},
		DockerImage: server.DockerImage,
		Bid:         server.BID,
	}

	return connect.NewResponse(&serverProto), nil
}
