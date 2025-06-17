package auth

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/global"
	"panelium/backend/internal/errors"
	"panelium/backend/model"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) Login(
	ctx context.Context,
	req *connect.Request[proto_gen_go.LoginRequest],
) (*connect.Response[proto_gen_go.LoginResponse], error) {
	result := global.DB.First(&model.User{}, "username = ? OR email = ?", req.Msg.Username, req.Msg.Username)
	if result.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.UserNotFound)
	}
	user := &model.User{}
	if err := result.Scan(user); err.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, err.Error)
	}

	// TODO: check password with salt and pepper against the hash
	// TODO: check if user needs 2fa
	// TODO: generate JWT auth+refresh tokens

	return nil, nil
}

// TODO: decide how to store pepper and implement it (at least 112 bits based on NIST recommendations - salt should be at least 128 bits)
