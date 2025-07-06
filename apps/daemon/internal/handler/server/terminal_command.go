package server

import (
	"connectrpc.com/connect"
	"context"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
)

func (s *ServerServiceHandler) TerminalCommand(
	ctx context.Context,
	req *connect.Request[proto_gen_go.IDMessage],
) (*connect.Response[proto_gen_go.Empty], error) {
	err := security.CheckServerAccess(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	err = server.TerminalCommand(req.Msg.Id, req.Msg.Text)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return nil, nil
}
