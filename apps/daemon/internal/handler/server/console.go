package server

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
)

func (s *ServerServiceHandler) Console(
	ctx context.Context,
	stm *connect.BidiStream[proto_gen_go.StreamIDMessage, proto_gen_go.SimpleMessage],
) error {
	firstMsg, err := stm.Receive()
	if err != nil {
		return err
	}
	if firstMsg.Id == nil || *firstMsg.Id == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("invalid server ID"))
	}

	err = security.CheckServerAccess(ctx, *firstMsg.Id)
	if err != nil {
		return connect.NewError(connect.CodePermissionDenied, err)
	}

	err = server.Console(*firstMsg.Id, stm)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	return nil
}
