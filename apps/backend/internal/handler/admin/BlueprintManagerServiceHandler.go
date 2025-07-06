package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/proto_gen_go/backend/admin"
)

type BlueprintManagerServiceHandler struct{}

func NewBlueprintManagerServiceHandler() *BlueprintManagerServiceHandler {
	return &BlueprintManagerServiceHandler{}
}

func (h *BlueprintManagerServiceHandler) GetBlueprints(ctx context.Context, req *connect.Request[admin.GetBlueprintsRequest]) (*connect.Response[admin.GetBlueprintsResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetBlueprintsResponse{}), nil
}

func (h *BlueprintManagerServiceHandler) GetBlueprint(ctx context.Context, req *connect.Request[admin.GetBlueprintRequest]) (*connect.Response[admin.GetBlueprintResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetBlueprintResponse{}), nil
}

func (h *BlueprintManagerServiceHandler) CreateBlueprint(ctx context.Context, req *connect.Request[admin.CreateBlueprintRequest]) (*connect.Response[admin.CreateBlueprintResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.CreateBlueprintResponse{}), nil
}

func (h *BlueprintManagerServiceHandler) UpdateBlueprint(ctx context.Context, req *connect.Request[admin.UpdateBlueprintRequest]) (*connect.Response[admin.UpdateBlueprintResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.UpdateBlueprintResponse{}), nil
}

func (h *BlueprintManagerServiceHandler) DeleteBlueprint(ctx context.Context, req *connect.Request[admin.DeleteBlueprintRequest]) (*connect.Response[admin.DeleteBlueprintResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.DeleteBlueprintResponse{}), nil
}
