package middleware

import (
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"net/http"
	"panelium/backend/internal/config"
)

func WithCORS(h http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{config.ConfigInstance.GetDashboardHost()},
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   connectcors.AllowedHeaders(),
		ExposedHeaders:   connectcors.ExposedHeaders(),
	})
	return corsMiddleware.Handler(h)
}
