package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/proto_gen_go/backend/admin"
)

type NodeAllocationManagerServiceHandler struct{}

func NewNodeAllocationManagerServiceHandler() *NodeAllocationManagerServiceHandler {
	return &NodeAllocationManagerServiceHandler{}
}

func (h *NodeAllocationManagerServiceHandler) GetNodeAllocations(ctx context.Context, req *connect.Request[admin.GetNodeAllocationsRequest]) (*connect.Response[admin.GetNodeAllocationsResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetNodeAllocationsResponse{}), nil
}

func (h *NodeAllocationManagerServiceHandler) GetNodeAllocation(ctx context.Context, req *connect.Request[admin.GetNodeAllocationRequest]) (*connect.Response[admin.GetNodeAllocationResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetNodeAllocationResponse{}), nil
}

func (h *NodeAllocationManagerServiceHandler) CreateNodeAllocation(ctx context.Context, req *connect.Request[admin.CreateNodeAllocationRequest]) (*connect.Response[admin.CreateNodeAllocationResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.CreateNodeAllocationResponse{}), nil
}

func (h *NodeAllocationManagerServiceHandler) UpdateNodeAllocation(ctx context.Context, req *connect.Request[admin.UpdateNodeAllocationRequest]) (*connect.Response[admin.UpdateNodeAllocationResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.UpdateNodeAllocationResponse{}), nil
}

func (h *NodeAllocationManagerServiceHandler) DeleteNodeAllocation(ctx context.Context, req *connect.Request[admin.DeleteNodeAllocationRequest]) (*connect.Response[admin.DeleteNodeAllocationResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.DeleteNodeAllocationResponse{}), nil
}
