package middleware

import (
	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"context"
	"net/http"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/common/errors"
	"panelium/common/jwt"
)

var AuthenticationMiddleware = authn.NewMiddleware(authentication)

var errorInvalidCredentials = connect.NewError(401, errors.InvalidCredentials)

type SessionInfo struct {
	SessionID string
	UserID    string
}

func authentication(_ context.Context, req *http.Request) (any, error) {
	accessTokenCookie, err := req.Cookie("access_jwt")
	if err != nil || accessTokenCookie.Value == "" {
		return nil, errorInvalidCredentials
	}

	claims, err := jwt.VerifyJWT(accessTokenCookie.Value, &config.JWTPrivateKeyInstance.PublicKey, jwt.BackendIssuer, jwt.AccessTokenType)
	if err != nil {
		return nil, errorInvalidCredentials
	}

	result := db.Instance().Model(&model.UserSession{}).First(&model.UserSession{}, "session_id = ? AND user_id = ?", claims.Audience, *claims.Subject)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeInternal, errors.SessionNotFound)
	}
	session := &model.UserSession{}
	if err := result.Scan(session); err.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, err.Error)
	}

	if session.AccessJTI != claims.JTI {
		// possible replay attack - delete the session to log out the user
		db.Instance().Model(&model.UserSession{}).Where("session_id = ?", claims.Audience).Delete(&model.UserSession{})
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.InvalidCredentials)
	}

	return SessionInfo{
		SessionID: claims.Audience,
		UserID:    *claims.Subject,
	}, nil
}
