package handler

import (
	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"panelium/daemon/internal/handler/server"
	"panelium/daemon/internal/handler/servers"
	"panelium/daemon/internal/middleware"
	"panelium/proto_gen_go/daemon/daemonconnect"
)

func Handle(host string) error {
	backendAuthInterceptors := connect.WithInterceptors(middleware.NewBackendAuthInterceptor())
	userAuthInterceptors := connect.WithInterceptors(middleware.NewUserAuthInterceptor())

	mux := http.NewServeMux()
	mux.Handle(daemonconnect.NewServersServiceHandler(&servers.ServersServiceHandler{}, backendAuthInterceptors))
	mux.Handle(daemonconnect.NewServerServiceHandler(&server.ServerServiceHandler{}, userAuthInterceptors))

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
