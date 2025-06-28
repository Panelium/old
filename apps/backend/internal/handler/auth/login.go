package auth

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/backend/internal/security"
	"panelium/backend/internal/security/session"
	"panelium/common/errors"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) Login(
	ctx context.Context,
	req *connect.Request[proto_gen_go.LoginRequest],
) (*connect.Response[proto_gen_go.LoginResponse], error) {
	user := &model.User{}
	tx := db.Instance().First(user, "username = ? OR email = ?", req.Msg.Username, req.Msg.Username)
	if tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.UserNotFound)
	}

	passwordValid := security.VerifyPassword(req.Msg.Password, user.PasswordSalt, user.PasswordHash)
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

	_, refreshToken, accessToken, err := session.CreateSession(user.UID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.SessionCreationFailed)
	}

	res := connect.NewResponse(&proto_gen_go.LoginResponse{
		Success:     true,
		RequiresMfa: false,
	})

	noop(refreshToken, accessToken) // TODO: remove this, just so go doesn't complain about unused variables

	/* TODO: COOKIES
	refresh_jwt: refreshToken,
	access_jwt: accessToken,
	*/

	return res, nil
}
