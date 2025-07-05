package auth

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/backend/internal/security"
	"panelium/backend/internal/security/cookies"
	"panelium/backend/internal/security/session"
	"panelium/common/errors"
	"panelium/proto_gen_go/backend"
)

// TODO: need to add rate limiting

func (s *AuthServiceHandler) Login(
	ctx context.Context,
	req *connect.Request[backend.LoginRequest],
) (*connect.Response[backend.LoginResponse], error) {
	user := &model.User{}
	tx := db.Instance().First(user, "username = ? OR email = ?", req.Msg.Username, req.Msg.Username)
	if tx.RowsAffected == 0 || tx.Error != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.UserNotFound)
	}

	passwordValid := security.VerifyPassword(req.Msg.Password, user.PasswordSalt, user.PasswordHash)
	if !passwordValid {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.InvalidCredentials)
	}

	if user.MFANeeded {
		//TODO: handle MFA
		//TODO: generate MFA session token

		res := connect.NewResponse(&backend.LoginResponse{
			RequiresMfa: true,
		})

		return res, nil
	}

	_, refreshToken, accessToken, refreshTokenExpiration, accessTokenExpiration, err := session.CreateSession(user.UID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.SessionCreationFailed)
	}

	res := connect.NewResponse(&backend.LoginResponse{
		Success:     true,
		RequiresMfa: false,
	})

	cookies.SetJWTCookie(res.Header(), "refresh_jwt", refreshToken, refreshTokenExpiration)
	cookies.SetJWTCookie(res.Header(), "access_jwt", accessToken, accessTokenExpiration)

	return res, nil
}
