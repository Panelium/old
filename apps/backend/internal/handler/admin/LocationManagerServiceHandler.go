package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/proto_gen_go/backend/admin"
)

type LocationManagerServiceHandler struct{}

func NewLocationManagerServiceHandler() *LocationManagerServiceHandler {
	return &LocationManagerServiceHandler{}
}

func (h *LocationManagerServiceHandler) GetLocations(ctx context.Context, req *connect.Request[admin.GetLocationsRequest]) (*connect.Response[admin.GetLocationsResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetLocationsResponse{}), nil
}

func (h *LocationManagerServiceHandler) GetLocation(ctx context.Context, req *connect.Request[admin.GetLocationRequest]) (*connect.Response[admin.GetLocationResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetLocationResponse{}), nil
}

func (h *LocationManagerServiceHandler) CreateLocation(ctx context.Context, req *connect.Request[admin.CreateLocationRequest]) (*connect.Response[admin.CreateLocationResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.CreateLocationResponse{}), nil
}

func (h *LocationManagerServiceHandler) UpdateLocation(ctx context.Context, req *connect.Request[admin.UpdateLocationRequest]) (*connect.Response[admin.UpdateLocationResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.UpdateLocationResponse{}), nil
}

func (h *LocationManagerServiceHandler) DeleteLocation(ctx context.Context, req *connect.Request[admin.DeleteLocationRequest]) (*connect.Response[admin.DeleteLocationResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.DeleteLocationResponse{}), nil
}
