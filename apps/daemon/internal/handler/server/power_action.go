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

func (s *ServerServiceHandler) PowerAction(
	ctx context.Context,
	req *connect.Request[proto_gen_go.PowerActionMessage],
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

	var err error

	switch req.Msg.Action {
	case proto_gen_go.PowerAction_POWER_ACTION_START:
		err = server.Start(srv)
	case proto_gen_go.PowerAction_POWER_ACTION_RESTART:
		err = server.Restart(srv)
	case proto_gen_go.PowerAction_POWER_ACTION_STOP:
		err = server.Stop(srv, false)
	case proto_gen_go.PowerAction_POWER_ACTION_KILL:
		err = server.Stop(srv, true)

	case proto_gen_go.PowerAction_POWER_ACTION_UNSPECIFIED:
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid power action"))
	}

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to perform power action"))
	}

	res := connect.NewResponse(&proto_gen_go.SuccessMessage{
		Success: true,
	})

	return res, nil
}
