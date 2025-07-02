package server

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *ServerServiceHandler) DeleteServer(
	ctx context.Context,
	req *connect.Request[proto_gen_go.DeleteServerRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	return nil, nil
}
