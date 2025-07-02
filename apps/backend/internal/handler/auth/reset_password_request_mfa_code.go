package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/proto_gen_go"
)

func (s *AuthServiceHandler) ResetPasswordRequestMFACode(
	ctx context.Context,
	req *connect.Request[proto_gen_go.RequestMFACodeRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
