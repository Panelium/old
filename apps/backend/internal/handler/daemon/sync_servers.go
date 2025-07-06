package daemon

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/backend/internal/middleware"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *DaemonServiceHandler) SyncServers(
	ctx context.Context,
	req *connect.Request[proto_gen_go.Empty],
	stm *connect.ServerStream[backend.Server],
) error {
	daemonInfoData := ctx.Value("panelium_daemon_info")
	daemonInfo, ok := daemonInfoData.(*middleware.DaemonInfo)
	if !ok || daemonInfo == nil || daemonInfo.NID == "" {
		return connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
	}

	return nil
}
