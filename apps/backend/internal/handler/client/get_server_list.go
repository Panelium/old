package client

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"panelium/backend/internal/db"
	"panelium/backend/internal/middleware"
	"panelium/backend/internal/model"
	"panelium/common/errors"
	"panelium/common/util"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *ClientServiceHandler) GetServerList(ctx context.Context, req *connect.Request[proto_gen_go.Empty]) (*connect.Response[backend.ServerList], error) {
	sessionInfoData := ctx.Value("panelium_session_info")
	sessionInfo, ok := sessionInfoData.(*middleware.SessionInfo)
	if !ok || sessionInfo == nil || sessionInfo.SessionID == "" || sessionInfo.UserID == "" {
		return nil, errors.ConnectInvalidCredentials
	}

	var user *model.User
	tx := db.Instance().First(&user, "uid = ?", sessionInfo.UserID)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.UserNotFound)
	}

	var servers []model.Server
	tx = db.Instance().Where("owner_id = ?", user.ID).Find(&servers)
	if tx.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to fetch servers"))
	}

	var accessibleServerUsers []model.ServerUser
	tx = db.Instance().Where("user_id = ?", user.ID).Find(&accessibleServerUsers)
	if tx.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to fetch accessible servers"))
	}

	var accessibleServers []model.Server
	for _, serverUser := range accessibleServerUsers {
		var server model.Server
		tx = db.Instance().First(&server, "id = ?", serverUser.ServerID)
		if tx.Error != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to fetch accessible server"))
		}
		accessibleServers = append(accessibleServers, server)
	}

	servers = append(servers, accessibleServers...)

	serverInfos := make([]*backend.ServerInfo, 0, len(servers))
	for _, server := range servers {
		serverInfo := &backend.ServerInfo{
			Sid:         server.SID,
			Name:        server.Name,
			Description: server.Description,
			Software:    server.Blueprint.Name,
			MainAllocation: util.IfElse(len(server.Allocations) > 0,
				&proto_gen_go.IPAllocation{
					Ip:   server.Allocations[0].IP,
					Port: uint32(server.Allocations[0].Port),
				},
				nil),
			DaemonHost: util.IfElse(server.Node.HTTPS, "https://", "http://") + server.Node.FQDN + ":" + fmt.Sprint(server.Node.DaemonPort),
		}
		serverInfos = append(serverInfos, serverInfo)
	}

	res := &backend.ServerList{
		Servers: serverInfos,
	}

	return connect.NewResponse(res), nil
}
