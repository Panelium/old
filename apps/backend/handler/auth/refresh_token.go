package auth

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) RefreshToken(
	context.Context,
	*connect.Request[proto_gen_go.RefreshTokenRequest],
) (*connect.Response[proto_gen_go.RefreshTokenResponse], error) {
	return nil, nil
}
