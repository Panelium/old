package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/proto_gen_go"
)

func (s *AuthServiceHandler) ResetPassword(
	ctx context.Context,
	req *connect.Request[proto_gen_go.ResetPasswordRequest],
) (*connect.Response[proto_gen_go.ResetPasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
