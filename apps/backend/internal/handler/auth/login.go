package auth

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/errors"
	"panelium/backend/internal/global"
	"panelium/backend/internal/model"
	"panelium/backend/internal/security"
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

	passwordValid := security.VerifyPassword(req.Msg.Password, user.PasswordSalt, global.Pepper, user.PasswordHash)
	if !passwordValid {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.InvalidCredentials)
	}

	if user.MFANeeded {
		//TODO: handle MFA
		//TODO: generate MFA session token

		res := connect.NewResponse(&proto_gen_go.LoginResponse{
			RequiresMfa: true,
		})

		return res, nil
	}

	// TODO: generate JWT auth+refresh tokens

	res := connect.NewResponse(&proto_gen_go.LoginResponse{
		RequiresMfa: false,
	})

	return res, nil
}
