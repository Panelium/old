package auth

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) Register(
	ctx context.Context,
	req *connect.Request[proto_gen_go.RegisterRequest],
) (*connect.Response[proto_gen_go.RegisterResponse], error) {
	return nil, nil
}
