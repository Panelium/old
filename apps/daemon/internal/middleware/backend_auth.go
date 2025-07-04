package middleware

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/common/jwt"
	"panelium/daemon/internal/config"
)

func NewBackendAuthInterceptor() connect.UnaryInterceptorFunc {
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

			claims, err := jwt.VerifyJWT(nodeToken, &config.JWTPrivateKeyInstance.PublicKey, jwt.DaemonIssuer, jwt.NodeTokenType)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
			}

			if claims.JTI != config.SecretsInstance.NodeJTI {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
			}

			return next(ctx, req)
		}
	}
	return interceptor
}
