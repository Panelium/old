package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/proto_gen_go/backend/admin"
)

type NodeManagerServiceHandler struct{}

func NewNodeManagerServiceHandler() *NodeManagerServiceHandler {
	return &NodeManagerServiceHandler{}
}

func (h *NodeManagerServiceHandler) GetNodes(ctx context.Context, req *connect.Request[admin.GetNodesRequest]) (*connect.Response[admin.GetNodesResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetNodesResponse{}), nil
}

func (h *NodeManagerServiceHandler) GetNode(ctx context.Context, req *connect.Request[admin.GetNodeRequest]) (*connect.Response[admin.GetNodeResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetNodeResponse{}), nil
}

func (h *NodeManagerServiceHandler) CreateNode(ctx context.Context, req *connect.Request[admin.CreateNodeRequest]) (*connect.Response[admin.CreateNodeResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.CreateNodeResponse{}), nil
}

func (h *NodeManagerServiceHandler) UpdateNode(ctx context.Context, req *connect.Request[admin.UpdateNodeRequest]) (*connect.Response[admin.UpdateNodeResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.UpdateNodeResponse{}), nil
}

func (h *NodeManagerServiceHandler) DeleteNode(ctx context.Context, req *connect.Request[admin.DeleteNodeRequest]) (*connect.Response[admin.DeleteNodeResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.DeleteNodeResponse{}), nil
}
