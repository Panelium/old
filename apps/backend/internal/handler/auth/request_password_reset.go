package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) RequestPasswordReset(
	ctx context.Context,
	req *connect.Request[proto_gen_go.RequestPasswordResetRequest],
) (*connect.Response[proto_gen_go.RequestPasswordResetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
}
