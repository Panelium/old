package auth

import (
	"connectrpc.com/connect"
	"context"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) VerifyMFA(
	ctx context.Context,
	req *connect.Request[proto_gen_go.VerifyMFARequest],
) (*connect.Response[proto_gen_go.VerifyMFAResponse], error) {
	return nil, nil
}
