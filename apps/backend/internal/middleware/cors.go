package middleware

import (
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"net/http"
)

func WithCORS(h http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // TODO: this needs to actually be configured properly
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return corsMiddleware.Handler(h)
}
