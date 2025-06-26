package middleware

import (
	"connectrpc.com/authn"
	"context"
	"net/http"
)

var AuthenticationMiddleware = authn.NewMiddleware(authentication)

func authentication(_ context.Context, req *http.Request) (any, error) {
	accessTokenCookie, err := req.Cookie("access_jwt")
	if err != nil {
		return nil, authn.Errorf("no access token")
	}

	//do stuff, verify jwt, etc

	return accessTokenCookie.Value, nil
}
