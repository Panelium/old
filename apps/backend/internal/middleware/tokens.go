package middleware

import (
	"connectrpc.com/authn"
	"context"
	"net/http"
	"strings"
)

var TokensMiddleware = authn.NewMiddleware(tokens)

type Tokens map[string]string

// tokens extracts JWT tokens from cookies and stores them in the context as a map
func tokens(ctx context.Context, req *http.Request) (any, error) {
	cookies := req.Cookies()

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

	context.WithValue(ctx, "panelium_tokens", tokenMap)

	return nil, nil // we don't want authn to modify the context
}
