package handler

import (
	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"panelium/backend/internal/handler/auth"
	"panelium/backend/internal/handler/client"
	"panelium/backend/internal/handler/daemon"
	"panelium/backend/internal/middleware"
	"panelium/proto_gen_go/backend/backendconnect"
)

func Handle(host string) error {
	daemonAuthInterceptors := connect.WithInterceptors(
		middleware.NewDaemonAuthInterceptor(),
	)

	userAuthInterceptors := connect.WithInterceptors(
		middleware.NewTokensInterceptor(),
		middleware.NewUserAuthInterceptor(),
	)

	mux := http.NewServeMux()
	mux.Handle(backendconnect.NewDaemonConnectionServiceHandler(&daemon.DaemonServiceHandler{}, daemonAuthInterceptors))

	mux.Handle(backendconnect.NewAuthServiceHandler(&auth.AuthServiceHandler{}, userAuthInterceptors))
	mux.Handle(backendconnect.NewClientServiceHandler(&client.ClientServiceHandler{}, userAuthInterceptors))

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
