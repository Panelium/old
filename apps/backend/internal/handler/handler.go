package handler

import (
	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"panelium/backend/internal/handler/auth"
	"panelium/backend/internal/middleware"
	"panelium/proto_gen_go/proto_gen_goconnect"
)

func Handle(host string) error {
	authInterceptors := connect.WithInterceptors(middleware.NewTokensInterceptor(), middleware.NewAuthInterceptor())

	mux := http.NewServeMux()
	mux.Handle(proto_gen_goconnect.NewAuthServiceHandler(&auth.AuthServiceHandler{}, authInterceptors))

	handler := h2c.NewHandler(mux, &http2.Server{})
	corsHandler := middleware.WithCORS(handler)
	err := http.ListenAndServe(
		host,
		corsHandler,
	)
	if err != nil {
		return err
	}
	return nil
}
