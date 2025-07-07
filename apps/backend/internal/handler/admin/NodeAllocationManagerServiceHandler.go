package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend/admin"
)

type NodeAllocationManagerServiceHandler struct{}

func NewNodeAllocationManagerServiceHandler() *NodeAllocationManagerServiceHandler {
	return &NodeAllocationManagerServiceHandler{}
}

func (h *NodeAllocationManagerServiceHandler) GetNodeAllocations(ctx context.Context, req *connect.Request[admin.GetNodeAllocationsRequest]) (*connect.Response[admin.GetNodeAllocationsResponse], error) {
	dbInst := db.Instance()
	var allocations []model.NodeAllocation
	var count int64
	page := req.Msg.Pagination.GetPage()
	pageSize := req.Msg.Pagination.GetPageSize()
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	query := dbInst.Model(&model.NodeAllocation{}).Preload("Node").Preload("Server")
	if req.Msg.Nid != nil && *req.Msg.Nid != "" {
		query = query.Joins("Node").Where("nodes.nid = ?", *req.Msg.Nid)
	}
	if req.Msg.Sid != nil && *req.Msg.Sid != "" {
		query = query.Joins("Server").Where("servers.sid = ?", *req.Msg.Sid)
	}
	query.Count(&count)
	query.Order("id desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&allocations)
	resp := &admin.GetNodeAllocationsResponse{
		NodeAllocations: make([]*admin.NodeAllocation, 0, len(allocations)),
		Pagination: &proto_gen_go.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    (*uint32)(nil),
		},
	}
	if count > 0 {
		total := uint32(count)
		resp.Pagination.Total = &total
	}
	for _, na := range allocations {
		resp.NodeAllocations = append(resp.NodeAllocations, NodeAllocationModelToProto(&na))
	}
	return connect.NewResponse(resp), nil
}

func (h *NodeAllocationManagerServiceHandler) GetNodeAllocation(ctx context.Context, req *connect.Request[admin.GetNodeAllocationRequest]) (*connect.Response[admin.GetNodeAllocationResponse], error) {
	dbInst := db.Instance()
	var allocation model.NodeAllocation
	if err := dbInst.Preload("Node").Preload("Server").Where("id = ?", req.Msg.Id).First(&allocation).Error; err != nil {
		return nil, err
	}
	resp := &admin.GetNodeAllocationResponse{
		NodeAllocation: NodeAllocationModelToProto(&allocation),
	}
	return connect.NewResponse(resp), nil
}

func (h *NodeAllocationManagerServiceHandler) CreateNodeAllocation(ctx context.Context, req *connect.Request[admin.CreateNodeAllocationRequest]) (*connect.Response[admin.CreateNodeAllocationResponse], error) {
	if req.Msg.NodeAllocation.Id != 0 {
		req.Msg.NodeAllocation.Id = 0
	}
	dbInst := db.Instance()
	na := NodeAllocationProtoToModel(req.Msg.NodeAllocation)
	// Find node by NID to set NodeID
	var node model.Node
	if err := dbInst.Where("nid = ?", req.Msg.NodeAllocation.Nid).First(&node).Error; err != nil {
		return nil, err
	}
	na.NodeID = node.ID
	// Optionally set ServerID if SID is provided
	if req.Msg.NodeAllocation.Sid != nil && *req.Msg.NodeAllocation.Sid != "" {
		var server model.Server
		if err := dbInst.Where("sid = ?", *req.Msg.NodeAllocation.Sid).First(&server).Error; err != nil {
			return nil, err
		}
		na.ServerID = &server.ID
	}
	if err := dbInst.Create(na).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.CreateNodeAllocationResponse{Success: true}), nil
}

func (h *NodeAllocationManagerServiceHandler) UpdateNodeAllocation(ctx context.Context, req *connect.Request[admin.UpdateNodeAllocationRequest]) (*connect.Response[admin.UpdateNodeAllocationResponse], error) {
	dbInst := db.Instance()
	na := NodeAllocationProtoToModel(req.Msg.NodeAllocation)
	// Find node by NID to set NodeID
	var node model.Node
	if err := dbInst.Where("nid = ?", req.Msg.NodeAllocation.Nid).First(&node).Error; err != nil {
		return nil, err
	}
	na.NodeID = node.ID
	// Optionally set ServerID if SID is provided
	if req.Msg.NodeAllocation.Sid != nil && *req.Msg.NodeAllocation.Sid != "" {
		var server model.Server
		if err := dbInst.Where("sid = ?", *req.Msg.NodeAllocation.Sid).First(&server).Error; err != nil {
			return nil, err
		}
		na.ServerID = &server.ID
	}
	if err := dbInst.Model(&model.NodeAllocation{}).Where("id = ?", req.Msg.NodeAllocation.Id).Updates(na).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.UpdateNodeAllocationResponse{Success: true}), nil
}

func (h *NodeAllocationManagerServiceHandler) DeleteNodeAllocation(ctx context.Context, req *connect.Request[admin.DeleteNodeAllocationRequest]) (*connect.Response[admin.DeleteNodeAllocationResponse], error) {
	dbInst := db.Instance()
	if err := dbInst.Where("id = ?", req.Msg.Id).Delete(&model.NodeAllocation{}).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.DeleteNodeAllocationResponse{Success: true}), nil
}
