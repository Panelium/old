package auth

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) RequestMFACode(
	ctx context.Context,
	req *connect.Request[proto_gen_go.RequestMFACodeRequest],
) (*connect.Response[proto_gen_go.RequestMFACodeResponse], error) {
	return nil, nil
}
