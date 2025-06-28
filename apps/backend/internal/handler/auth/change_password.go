package auth

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func ChangePassword(ctx context.Context, req *connect.Request[proto_gen_go.ChangePasswordRequest]) (*connect.Response[proto_gen_go.ChangePasswordResponse], error) {
	return nil, nil
}
