package server

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
)

func (s *ServerServiceHandler) Console(
	ctx context.Context,
	stm *connect.BidiStream[proto_gen_go.SimpleMessage, proto_gen_go.SimpleMessage],
) error {
	serverId := ctx.Value("server_id").(string)
	if serverId == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("server ID is required"))
	}

	err := server.Console(serverId, stm)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	return nil
}
