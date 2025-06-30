package middleware

import (
	"connectrpc.com/connect"
	"context"
	"net/http"
	"strings"
)

type Tokens map[string]string

func NewTokensInterceptor() connect.UnaryInterceptorFunc {
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
				return next(ctx, req)
			}

			cookies, err := http.ParseCookie(cookieString)
			if err != nil {
				return next(ctx, req)
			}

			tokenMap := Tokens{}

			for _, cookie := range cookies {
				if !strings.HasSuffix(cookie.Name, "_jwt") {
					continue
				}
				if cookie.Value == "" {
					continue
				}

				if !cookie.Secure {
					continue
				}
				if !cookie.HttpOnly {
					continue
				}
				if cookie.SameSite != http.SameSiteStrictMode {
					continue
				}

				tokenMap[cookie.Name] = cookie.Value
			}

			ctx = context.WithValue(ctx, "panelium_tokens", tokenMap)

			return next(ctx, req)
		}
	}
	return interceptor
}
