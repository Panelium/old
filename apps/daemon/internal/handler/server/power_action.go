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

func (s *ServerServiceHandler) PowerAction(
	ctx context.Context,
	req *connect.Request[daemon.PowerActionMessage],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	err := security.CheckServerAccess(ctx, req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	}

	var srv *model.Server
	tx := db.Instance().First(&srv, "sid = ?", req.Msg.ServerId)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("server not found"))
	}

	switch req.Msg.Action {
	case daemon.PowerAction_POWER_ACTION_START:
		err = server.Start(srv.SID)
	case daemon.PowerAction_POWER_ACTION_RESTART:
		err = server.Restart(srv.SID)
	case daemon.PowerAction_POWER_ACTION_STOP:
		err = server.Stop(srv.SID, false)
	case daemon.PowerAction_POWER_ACTION_KILL:
		err = server.Stop(srv.SID, true)

	case daemon.PowerAction_POWER_ACTION_UNSPECIFIED:
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
