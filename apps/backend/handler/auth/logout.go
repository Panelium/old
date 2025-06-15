package auth

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) Logout(
	context.Context,
	*connect.Request[proto_gen_go.LogoutRequest],
) (*connect.Response[proto_gen_go.LogoutResponse], error) {
	return nil, nil
}
