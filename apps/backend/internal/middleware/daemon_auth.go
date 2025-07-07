package middleware

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/common/jwt"
	"slices"
)

var daemonAuthIgnoredProcedures = []string{}

type DaemonInfo struct {
	NID string
}

func NewDaemonAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}
			if slices.Contains(daemonAuthIgnoredProcedures, req.Spec().Procedure) {
				return next(ctx, req)
			}

			nodeToken := req.Header().Get("Authorization")
			if nodeToken == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing node token"))
			}

			claims, err := jwt.VerifyJWT(nodeToken, &config.JWTPrivateKeyInstance.PublicKey, jwt.BackendIssuer, jwt.BackendTokenType)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
			}

			var node *model.Node
			tx := db.Instance().First(node, "backend_jti = ?", claims.JTI)
			if tx.Error != nil || tx.RowsAffected == 0 {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("node token not found"))
			}

			ctx = context.WithValue(ctx, "panelium_daemon_info", &DaemonInfo{
				NID: node.NID,
			})

			return next(ctx, req)
		}
	}
	return interceptor
}
