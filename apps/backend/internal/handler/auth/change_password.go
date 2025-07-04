package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/proto_gen_go/backend"
)

func (s *AuthServiceHandler) ChangePassword(
	ctx context.Context,
	req *connect.Request[backend.ChangePasswordRequest],
) (*connect.Response[backend.ChangePasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
