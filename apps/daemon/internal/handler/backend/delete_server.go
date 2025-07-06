package backend

import (
	"connectrpc.com/connect"
	"context"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go"
)

func (s *BackendServiceHandler) DeleteServer(
	ctx context.Context,
	req *connect.Request[proto_gen_go.SimpleIDMessage],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	err := server.DeleteServer(req.Msg.Id, true) // TODO: maybe don't always force delete?
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&proto_gen_go.SuccessMessage{
		Success: true,
	})

	return res, nil
}
