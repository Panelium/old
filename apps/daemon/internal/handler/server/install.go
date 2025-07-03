package server

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
)

func (s *ServerServiceHandler) Install(
	ctx context.Context,
	req *connect.Request[proto_gen_go.Empty],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	serverId := ctx.Value("server_id").(string)
	if serverId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("server ID is required"))
	}

	var srv *model.Server
	tx := db.Instance().First(&srv, "sid = ?", serverId)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("server not found"))
	}

	err := server.Install(srv)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to install server"))
	}

	res := connect.NewResponse(&proto_gen_go.SuccessMessage{
		Success: true,
	})

	return res, nil
}
