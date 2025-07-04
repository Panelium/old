package middleware

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"net/http"
	"panelium/common/jwt"
	"panelium/daemon/internal/config"
)

func NewUserAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			cookieString := req.Header().Get("Cookie")
			if cookieString == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
			}

			cookies, err := http.ParseCookie(cookieString)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
			}

			var cookie *http.Cookie

			for _, c := range cookies {
				if c.Name != "access_jwt" {
					continue
				}

				if c.Value == "" {
					continue
				}

				if !c.Secure {
					continue
				}
				if !c.HttpOnly {
					continue
				}
				if c.SameSite != http.SameSiteStrictMode {
					continue
				}

				cookie = c
			}

			if cookie == nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
			}

			accessToken := cookie.Value

			_, err = jwt.VerifyJWT(accessToken, config.BackendJWTPublicKeyInstance, jwt.BackendIssuer, jwt.AccessTokenType)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
			}

			// TODO: check server access??

			return next(ctx, req)
		}
	}
	return interceptor
}
