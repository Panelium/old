package middleware

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"log"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/common/jwt"
)

type DaemonInfo struct {
	NID string
}

func NewDaemonAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			log.Printf("HERE")
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			nodeToken := req.Header().Get("Authorization")
			if nodeToken == "" {
				log.Printf("missing node token in request header")
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing node token"))
			}

			claims, err := jwt.VerifyJWT(nodeToken, &config.JWTPrivateKeyInstance.PublicKey, jwt.BackendIssuer, jwt.BackendTokenType)
			if err != nil {
				log.Printf("failed to verify node token: %v", err)
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
			}

			var node *model.Node
			tx := db.Instance().First(&node, "backend_jti = ?", claims.JTI)
			if tx.Error != nil || tx.RowsAffected == 0 {
				log.Printf("error finding node with backend JTI %s: %v", claims.JTI, tx.Error)
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("node token not found"))
			}

			log.Printf("Node found: %s", node.NID)

			ctx = context.WithValue(ctx, "panelium_daemon_info", &DaemonInfo{
				NID: node.NID,
			})

			return next(ctx, req)
		}
	}
	return interceptor
}
