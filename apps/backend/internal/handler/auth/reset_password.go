package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/proto_gen_go/backend"
)

func (s *AuthServiceHandler) ResetPassword(
	ctx context.Context,
	req *connect.Request[backend.ResetPasswordRequest],
) (*connect.Response[backend.ResetPasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
