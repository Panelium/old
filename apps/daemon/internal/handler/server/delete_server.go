package server

import (
	"connectrpc.com/connect"
	"context"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/daemon"
)

func (s *ServerServiceHandler) DeleteServer(
	ctx context.Context,
	req *connect.Request[daemon.DeleteServerRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	err := server.DeleteServer(req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&proto_gen_go.SuccessMessage{
		Success: true,
	})

	return res, nil
}
