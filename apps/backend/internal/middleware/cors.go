package middleware

import (
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"net/http"
)

func WithCORS(h http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://127.0.0.1:5173", "http://127.0.0.1:9090"}, // TODO: this needs to actually be configured properly
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   connectcors.AllowedHeaders(),
		ExposedHeaders:   connectcors.ExposedHeaders(),
	})
	return corsMiddleware.Handler(h)
}
