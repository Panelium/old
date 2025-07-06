package daemon

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/backend/internal/middleware"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *DaemonServiceHandler) RegisterDaemon(
	ctx context.Context,
	req *connect.Request[backend.RegisterDaemonRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	daemonInfoData := ctx.Value("panelium_daemon_info")
	daemonInfo, ok := daemonInfoData.(*middleware.DaemonInfo)
	if !ok || daemonInfo == nil || daemonInfo.NID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
	}

	return nil, nil
}
