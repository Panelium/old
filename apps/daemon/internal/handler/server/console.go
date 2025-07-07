package server

import (
	"connectrpc.com/connect"
	"context"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
)

func (s *ServerServiceHandler) Console(
	ctx context.Context,
	req *connect.Request[proto_gen_go.SimpleIDMessage],
	stm *connect.ServerStream[proto_gen_go.SimpleMessage],
) error {
	err := security.CheckServerAccess(ctx, req.Msg.Id)
	if err != nil {
		return connect.NewError(connect.CodeFailedPrecondition, err)
	}

	err = server.Console(req.Msg.Id, stm)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	return nil
}
