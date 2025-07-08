package middleware

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/backend/internal/security/session"
	"panelium/common/errors"
	"panelium/common/jwt"
	"panelium/proto_gen_go/backend/backendconnect"
	"slices"
)

type SessionInfo struct {
	SessionID string
	UserID    string
}

var userAuthIgnoredProcedures = []string{
	backendconnect.AuthServiceRegisterProcedure,
	backendconnect.AuthServiceLoginProcedure,
	backendconnect.AuthServiceRequestMFACodeProcedure,
	backendconnect.AuthServiceRequestPasswordResetProcedure,
	backendconnect.AuthServiceResetPasswordProcedure,
	backendconnect.AuthServiceResetPasswordRequestMFACodeProcedure,
	backendconnect.AuthServiceResetPasswordVerifyMFAProcedure,
	backendconnect.AuthServiceVerifyMFAProcedure,
	backendconnect.AuthServiceRefreshTokenProcedure,
}

func NewUserAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}
			if slices.Contains(userAuthIgnoredProcedures, req.Spec().Procedure) {
				return next(ctx, req)
			}

			tokensData := ctx.Value("panelium_tokens")
			tokens, ok := tokensData.(Tokens)
			if !ok || tokens == nil || len(tokens) == 0 {
				return nil, errors.ConnectInvalidCredentials
			}

			accessToken, ok := tokens["access_jwt"]
			if ok != true {
				return nil, errors.ConnectInvalidCredentials // TODO: let the client know it should try refreshing the token (token is deleted when expired)
			}

			claims, err := jwt.VerifyJWT(accessToken, &config.JWTPrivateKeyInstance.PublicKey, jwt.BackendIssuer, jwt.AccessTokenType)
			if err != nil {
				return nil, errors.ConnectInvalidCredentials // TODO: let the client know it should try refreshing the token
			}

			userSession := &model.UserSession{}
			tx := db.Instance().Model(&model.UserSession{}).First(userSession, "session_id = ? AND user_id = ?", claims.Audience, *claims.Subject)
			if tx.Error != nil || tx.RowsAffected == 0 {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.SessionNotFound)
			}

			if userSession.AccessJTI != claims.JTI {
				// possible replay attack - delete the session to log out the users
				err := session.DeleteSession(userSession.SessionID)
				if err != nil {
					// TODO: log this error
					return nil, errors.ConnectInvalidCredentials
				}
				return nil, errors.ConnectInvalidCredentials
			}

			ctx = context.WithValue(ctx, "panelium_session_info", &SessionInfo{
				SessionID: *claims.Audience,
				UserID:    *claims.Subject,
			})

			return next(ctx, req)
		}
	}
	return interceptor
}
