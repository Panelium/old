package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *AuthServiceHandler) RequestPasswordReset(
	ctx context.Context,
	req *connect.Request[backend.RequestPasswordResetRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
