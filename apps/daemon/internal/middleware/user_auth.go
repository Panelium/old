package middleware

import (
	"connectrpc.com/connect"
	"context"
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
			// TODO
			//cookieString := req.Header().Get("Cookie")
			//if cookieString == "" {
			//	log.Printf("err %v: ", "cookie not found")
			//	return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
			//}
			//
			//cookies, err := http.ParseCookie(cookieString)
			//if err != nil {
			//	log.Printf("err %v: ", err)
			//	return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
			//}
			//
			//var cookie *http.Cookie
			//
			//for _, c := range cookies {
			//	if c.Name != "access_jwt" {
			//		continue
			//	}
			//
			//	if c.Value == "" {
			//		continue
			//	}
			//
			//	if !c.Secure {
			//		continue
			//	}
			//	if !c.HttpOnly {
			//		continue
			//	}
			//	if c.SameSite != http.SameSiteStrictMode {
			//		continue
			//	}
			//
			//	cookie = c
			//}
			//
			//if cookie == nil {
			//	log.Printf("err %v: ", err)
			//	return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
			//}
			//
			//accessToken := cookie.Value
			//
			//claims, err := jwt.VerifyJWT(accessToken, config.BackendJWTPublicKeyInstance, jwt.BackendIssuer, jwt.AccessTokenType)
			//if err != nil {
			//	log.Printf("err %v: ", err)
			//	return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid access token"))
			//}
			//
			//ctx = context.WithValue(ctx, "user_id", claims.Subject)

			return next(ctx, req)
		}
	}
	return interceptor
}
