package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/proto_gen_go/backend"
)

func (s *AuthServiceHandler) ChangePasswordVerifyMFA(
	ctx context.Context,
	req *connect.Request[backend.VerifyMFARequest],
) (*connect.Response[backend.VerifyMFAResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
