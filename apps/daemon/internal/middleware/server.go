package middleware

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
)

func NewServerInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			serverId := req.Header().Get("ServerID")
			if serverId == "" {
				return next(ctx, req)
			}

			tx := db.Instance().First(&model.Server{}, "server_id = ?", serverId)
			if tx.Error != nil || tx.RowsAffected == 0 {
				return nil, connect.NewError(connect.CodeNotFound, errors.New("server not found"))
			}

			ctx = context.WithValue(ctx, "server_id", serverId)

			return next(ctx, req)
		}
	}
	return interceptor
}
