package auth

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"gorm.io/gorm"
	"panelium/backend/internal/db"
	"panelium/backend/internal/global"
	"panelium/backend/internal/model"
	"panelium/backend/internal/rate_limit"
	"panelium/backend/internal/security"
	"panelium/backend/internal/security/cookies"
	"panelium/backend/internal/security/session"
	"panelium/common/id"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
	"time"
)

// TODO: need to add rate limiting

var registerLimiter = rate_limit.NewRateLimiter(5, time.Minute) // 5 requests/minute per IP

func (s *AuthServiceHandler) Register(
	ctx context.Context,
	req *connect.Request[backend.RegisterRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	if !registerLimiter.Allow(req.Peer().Addr) {
		return nil, connect.NewError(connect.CodeResourceExhausted, errors.New("too many register attempts, please try again later"))
	}

	err := global.ValidatorInstance().Var(req.Msg.Email, "required,email")
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email is missing or invalid"))
	}

	err = global.ValidatorInstance().Var(req.Msg.Username, "required,username")
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username is missing or invalid"))
	}

	err = global.ValidatorInstance().Var(req.Msg.Password, "required,min=16,max=384") // TODO: potentially make the password requirements configurable
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password is missing or invalid"))
	}

	// TODO: password blacklist? (rockyou or something)

	tx := db.Instance().Where("email = ? OR username = ?", req.Msg.Email, req.Msg.Username).First(&model.User{})
	if tx.Error == nil {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("email or username already exists"))
	}
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, connect.NewError(connect.CodeInternal, errors.New("database error"))
	}

	uid, err := id.New()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to create user"))
	}

	passwordHash, passwordSalt, err := security.HashPassword(req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to create user"))
	}

	user := &model.User{
		UID:          uid,
		Username:     req.Msg.Username,
		Email:        req.Msg.Email,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		MFANeeded:    false, // TODO: handle MFA in the future
	}

	if err := db.Instance().Create(user).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to create user"))
	}

	_, refreshToken, accessToken, refreshTokenExpiration, accessTokenExpiration, err := session.CreateSession(uid)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to create session"))
	}

	res := connect.NewResponse(&proto_gen_go.SuccessMessage{
		Success: true,
	})

	cookies.SetJWTCookie(res.Header(), "refresh_jwt", refreshToken, refreshTokenExpiration)
	cookies.SetJWTCookie(res.Header(), "access_jwt", accessToken, accessTokenExpiration)

	return res, nil
}
