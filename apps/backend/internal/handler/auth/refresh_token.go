package auth

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/middleware"
	"panelium/backend/internal/model"
	"panelium/backend/internal/security/cookies"
	"panelium/backend/internal/security/session"
	"panelium/common/errors"
	"panelium/common/jwt"
	"panelium/proto_gen_go"
)

func (s *AuthServiceHandler) RefreshToken(
	ctx context.Context,
	req *connect.Request[proto_gen_go.Empty],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	tokensData := ctx.Value("panelium_tokens")
	tokens, ok := tokensData.(middleware.Tokens)
	if !ok || tokens == nil || len(tokens) == 0 {
		return nil, errors.ConnectInvalidCredentials
	}

	refreshToken, ok := tokens["refresh_jwt"]
	if ok != true {
		return nil, errors.ConnectInvalidCredentials
	}

	claims, err := jwt.VerifyJWT(refreshToken, &config.JWTPrivateKeyInstance.PublicKey, jwt.BackendIssuer, jwt.RefreshTokenType)
	if err != nil {
		return nil, errors.ConnectInvalidCredentials
	}

	userSession := &model.UserSession{}
	tx := db.Instance().Model(&model.UserSession{}).First(userSession, "session_id = ? AND user_id = ?", claims.Audience, *claims.Subject)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.SessionNotFound)
	}

	if userSession.RefreshJTI != claims.JTI {
		err := session.DeleteSession(userSession.SessionID)
		if err != nil {
			// TODO: log this error
			return nil, errors.ConnectInvalidCredentials
		}
		return nil, errors.ConnectInvalidCredentials
	}

	newRefreshToken, newAccessToken, newRefreshTokenExpiration, newAccessTokenExpiration, err := session.RefreshSession(userSession.SessionID) // RefreshSession will also update the session in the database
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.SessionCreationFailed)
	}

	res := connect.NewResponse(&proto_gen_go.SuccessMessage{
		Success: true,
	})

	cookies.SetJWTCookie(res.Header(), "refresh_jwt", newRefreshToken, newRefreshTokenExpiration)
	cookies.SetJWTCookie(res.Header(), "access_jwt", newAccessToken, newAccessTokenExpiration)

	return res, nil
}
