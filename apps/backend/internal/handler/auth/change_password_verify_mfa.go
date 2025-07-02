package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/proto_gen_go"
)

func (s *AuthServiceHandler) ChangePasswordVerifyMFA(
	ctx context.Context,
	req *connect.Request[proto_gen_go.VerifyMFARequest],
) (*connect.Response[proto_gen_go.VerifyMFAResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
