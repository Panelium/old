package middleware

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/common/jwt"
)

func NewDaemonAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			nodeToken := req.Header().Get("Authorization")
			if nodeToken == "" {
				return next(ctx, req)
			}

			claims, err := jwt.VerifyJWT(nodeToken, &config.JWTPrivateKeyInstance.PublicKey, jwt.BackendIssuer, jwt.BackendTokenType)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
			}

			tx := db.Instance().First(&model.Node{}, "backend_jti = ?", claims.JTI)
			if tx.Error != nil || tx.RowsAffected == 0 {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
			}

			return next(ctx, req)
		}
	}
	return interceptor
}
