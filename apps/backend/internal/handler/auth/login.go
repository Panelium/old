package auth

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/errors"
	"panelium/backend/internal/global"
	"panelium/backend/internal/model"
	"panelium/backend/internal/security"
	"panelium/backend/internal/security/session"
	"panelium/common/jwt"
	proto_gen_go "panelium/proto-gen-go"
	"time"
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

	sessionId, refreshToken, err := session.CreateSession(user.UID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.SessionCreationFailed)
	}

	claims := jwt.Claims{
		IssuedAt:   time.Now().Unix(),
		NotBefore:  nil,
		Expiration: time.Now().Add(time.Minute * 5).Unix(), // TODO: this needs more thought, perhaps config?
		Subject:    &user.UID,
		Audience:   sessionId,
		Issuer:     "backend", // TODO: we might want to make this shorter
		TokenType:  "access",  // TODO: we might want to make this shorter
		JTI:        nil,       // TODO: add a JTI
	}
	accessToken, err := jwt.CreateJWT(claims, global.JWTSecret)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.TokenCreationFailed)
	}

	res := connect.NewResponse(&proto_gen_go.LoginResponse{
		RequiresMfa:  false,
		RefreshToken: &refreshToken, // TODO: this needs to be turned into a cookie
		AccessToken:  &accessToken,  // TODO: this needs to be turned into a cookie
	})

	return res, nil
}
