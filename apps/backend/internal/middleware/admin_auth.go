package middleware

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
)

func NewAdminAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			sessionInfoData := ctx.Value("panelium_session_info")
			sessionInfo, ok := sessionInfoData.(*SessionInfo)
			if !ok {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid session"))
			}
			if sessionInfo == nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("session not found"))
			}

			var user *model.User
			tx := db.Instance().First(&user, "uid = ?", sessionInfo.UserID)
			if tx.Error != nil || tx.RowsAffected == 0 {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user not found"))
			}

			if !user.Admin {
				return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user is not an admin"))
			}

			return next(ctx, req)
		}
	}
	return interceptor
}
