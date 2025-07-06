package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/proto_gen_go/backend/admin"
)

type ServerManagerServiceHandler struct{}

func NewServerManagerServiceHandler() *ServerManagerServiceHandler {
	return &ServerManagerServiceHandler{}
}

func (h *ServerManagerServiceHandler) GetServers(ctx context.Context, req *connect.Request[admin.GetServersRequest]) (*connect.Response[admin.GetServersResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetServersResponse{}), nil
}

func (h *ServerManagerServiceHandler) GetServer(ctx context.Context, req *connect.Request[admin.GetServerRequest]) (*connect.Response[admin.GetServerResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetServerResponse{}), nil
}

func (h *ServerManagerServiceHandler) CreateServer(ctx context.Context, req *connect.Request[admin.CreateServerRequest]) (*connect.Response[admin.CreateServerResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.CreateServerResponse{}), nil
}

func (h *ServerManagerServiceHandler) UpdateServer(ctx context.Context, req *connect.Request[admin.UpdateServerRequest]) (*connect.Response[admin.UpdateServerResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.UpdateServerResponse{}), nil
}

func (h *ServerManagerServiceHandler) DeleteServer(ctx context.Context, req *connect.Request[admin.DeleteServerRequest]) (*connect.Response[admin.DeleteServerResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.DeleteServerResponse{}), nil
}
