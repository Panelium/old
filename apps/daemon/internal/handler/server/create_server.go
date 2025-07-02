package server

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *ServerServiceHandler) CreateServer(context.Context, *connect.Request[proto_gen_go.CreateServerRequest]) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	return nil, nil
}
