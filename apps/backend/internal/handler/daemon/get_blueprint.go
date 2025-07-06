package daemon

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/backend/internal/middleware"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *DaemonServiceHandler) GetBlueprint(
	ctx context.Context,
	req *connect.Request[proto_gen_go.SimpleIDMessage],
) (*connect.Response[backend.Blueprint], error) {
	daemonInfoData := ctx.Value("panelium_daemon_info") // this isn't really necessary for this method rn
	daemonInfo, ok := daemonInfoData.(*middleware.DaemonInfo)
	if !ok || daemonInfo == nil || daemonInfo.NID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
	}

	return nil, nil
}
