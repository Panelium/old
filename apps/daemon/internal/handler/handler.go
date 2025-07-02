package handler

import (
	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"panelium/daemon/internal/handler/server"
	"panelium/daemon/internal/middleware"
	"panelium/proto-gen-go/proto_gen_goconnect"
)

func Handle(host string) error {
	authInterceptors := connect.WithInterceptors()

	mux := http.NewServeMux()
	mux.Handle(proto_gen_goconnect.NewServerServiceHandler(&server.ServerServiceHandler{}, authInterceptors))

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
