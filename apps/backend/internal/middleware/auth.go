package middleware

import (
	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"context"
	stdErrors "errors"
	"net/http"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/common/errors"
	"panelium/common/jwt"
	"panelium/proto-gen-go/proto_gen_goconnect"
	"slices"
)

var allowedProceduresAuth = []string{
	proto_gen_goconnect.AuthServiceRefreshTokenProcedure,
}

var AuthenticationMiddleware = authn.NewMiddleware(authentication)

type SessionInfo struct {
	SessionID string
	UserID    string
}

func authentication(ctx context.Context, req *http.Request) (any, error) {
	procedure, ok := authn.InferProcedure(req.URL)
	if !ok {
		return nil, stdErrors.New("huh? this should not happen")
	}
	if !slices.Contains(allowedProceduresAuth, procedure) {
		// If the procedure is not allowed, we don't need to process tokens.
		return nil, nil
	}

	tokens := ctx.Value("panelium_tokens").(Tokens)
	if tokens == nil || len(tokens) == 0 {
		return nil, errors.ConnectInvalidCredentials
	}

	accessToken, ok := tokens["access_jwt"]
	if ok != true {
		return nil, errors.ConnectInvalidCredentials
	}

	claims, err := jwt.VerifyJWT(accessToken, &config.JWTPrivateKeyInstance.PublicKey, jwt.BackendIssuer, jwt.AccessTokenType)
	if err != nil {
		return nil, errors.ConnectInvalidCredentials
	}

	session := &model.UserSession{}
	tx := db.Instance().Model(&model.UserSession{}).First(session, "session_id = ? AND user_id = ?", claims.Audience, *claims.Subject)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.SessionNotFound)
	}

	if session.AccessJTI != claims.JTI {
		// possible replay attack - delete the session to log out the user
		db.Instance().Model(&model.UserSession{}).Where("session_id = ?", claims.Audience).Delete(&model.UserSession{})
		return nil, errors.ConnectInvalidCredentials
	}

	return SessionInfo{
		SessionID: claims.Audience,
		UserID:    *claims.Subject,
	}, nil
	// we can get this data by calling authn.GetInfo(ctx)
}
