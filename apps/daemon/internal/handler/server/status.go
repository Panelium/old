package server

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/daemon"
)

func (s *ServerServiceHandler) Status(
	ctx context.Context,
	req *connect.Request[proto_gen_go.SimpleIDMessage],
) (*connect.Response[daemon.ServerStatus], error) {
	err := security.CheckServerAccess(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	var srv *model.Server
	tx := db.Instance().First(&srv, "sid = ?", req.Msg.Id)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("server not found"))
	}

	status := server.Status(srv)

	res := connect.NewResponse(status)

	return res, nil
}
