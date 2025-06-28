package auth

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func ResetPassword(ctx context.Context, req *connect.Request[proto_gen_go.ResetPasswordRequest]) (*connect.Response[proto_gen_go.ResetPasswordResponse], error) {
	return nil, nil
}
