package handler

import (
	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"panelium/backend/internal/handler/admin"
	"panelium/backend/internal/handler/auth"
	"panelium/backend/internal/handler/client"
	"panelium/backend/internal/handler/daemon"
	"panelium/backend/internal/middleware"
	"panelium/proto_gen_go/backend/admin/adminconnect"
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

	adminAuthInterceptors := connect.WithInterceptors(
		middleware.NewTokensInterceptor(),
		middleware.NewUserAuthInterceptor(),
		middleware.NewAdminAuthInterceptor(),
	)

	mux := http.NewServeMux()
	mux.Handle(backendconnect.NewDaemonServiceHandler(&daemon.DaemonServiceHandler{}, daemonAuthInterceptors))

	mux.Handle(backendconnect.NewAuthServiceHandler(&auth.AuthServiceHandler{}, userAuthInterceptors))
	mux.Handle(backendconnect.NewClientServiceHandler(&client.ClientServiceHandler{}, userAuthInterceptors))

	mux.Handle(adminconnect.NewUserManagerServiceHandler(admin.NewUserManagerServiceHandler(), adminAuthInterceptors))
	mux.Handle(adminconnect.NewBlueprintManagerServiceHandler(admin.NewBlueprintManagerServiceHandler(), adminAuthInterceptors))
	mux.Handle(adminconnect.NewLocationManagerServiceHandler(admin.NewLocationManagerServiceHandler(), adminAuthInterceptors))
	mux.Handle(adminconnect.NewNodeManagerServiceHandler(admin.NewNodeManagerServiceHandler(), adminAuthInterceptors))
	mux.Handle(adminconnect.NewNodeAllocationManagerServiceHandler(admin.NewNodeAllocationManagerServiceHandler(), adminAuthInterceptors))
	mux.Handle(adminconnect.NewServerManagerServiceHandler(admin.NewServerManagerServiceHandler(), adminAuthInterceptors))

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
