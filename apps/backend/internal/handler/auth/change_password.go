package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/proto_gen_go"
)

func (s *AuthServiceHandler) ChangePassword(
	ctx context.Context,
	req *connect.Request[proto_gen_go.ChangePasswordRequest],
) (*connect.Response[proto_gen_go.ChangePasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
