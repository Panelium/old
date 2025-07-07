package admin

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/backend/internal/security/session"
	"panelium/common/id"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend/admin"
	"panelium/proto_gen_go/backend/admin/adminconnect"
	"time"
)

type NodeManagerServiceHandler struct {
	adminconnect.NodeManagerServiceHandler
}

func NewNodeManagerServiceHandler() *NodeManagerServiceHandler {
	return &NodeManagerServiceHandler{}
}

func (h *NodeManagerServiceHandler) GetNodes(ctx context.Context, req *connect.Request[admin.GetNodesRequest]) (*connect.Response[admin.GetNodesResponse], error) {
	dbInst := db.Instance()
	var nodes []model.Node
	var count int64
	page := req.Msg.Pagination.GetPage()
	pageSize := req.Msg.Pagination.GetPageSize()
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	query := dbInst.Model(&model.Node{}).Preload("Location")
	if req.Msg.Lid != nil && *req.Msg.Lid != "" {
		query = query.Joins("Location").Where("locations.lid = ?", req.Msg.Lid)
	}
	query.Count(&count)
	query.Order("id desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&nodes)
	resp := &admin.GetNodesResponse{
		Nodes: make([]*admin.Node, 0, len(nodes)),
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
	for _, n := range nodes {
		resp.Nodes = append(resp.Nodes, NodeModelToProto(&n))
	}
	return connect.NewResponse(resp), nil
}

func (h *NodeManagerServiceHandler) GetNode(ctx context.Context, req *connect.Request[admin.GetNodeRequest]) (*connect.Response[admin.GetNodeResponse], error) {
	dbInst := db.Instance()
	var node model.Node
	if err := dbInst.Preload("Location").Where("nid = ?", req.Msg.Nid).First(&node).Error; err != nil {
		return nil, err
	}
	resp := &admin.GetNodeResponse{
		Node: NodeModelToProto(&node),
	}
	return connect.NewResponse(resp), nil
}

func (h *NodeManagerServiceHandler) CreateNode(ctx context.Context, req *connect.Request[admin.CreateNodeRequest]) (*connect.Response[admin.CreateNodeResponse], error) {
	var err error
	req.Msg.Node.Nid, err = id.New()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create node"))
	}

	dbInst := db.Instance()
	node := NodeProtoToModel(req.Msg.Node)
	// Find location by LID to set LocationID
	var location model.Location
	if err := dbInst.Where("lid = ?", req.Msg.Node.Lid).First(&location).Error; err != nil {
		return nil, err
	}
	node.LocationID = location.ID
	if err := dbInst.Create(node).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.CreateNodeResponse{Success: true}), nil
}

func (h *NodeManagerServiceHandler) UpdateNode(ctx context.Context, req *connect.Request[admin.UpdateNodeRequest]) (*connect.Response[admin.UpdateNodeResponse], error) {
	dbInst := db.Instance()
	node := NodeProtoToModel(req.Msg.Node)
	// Find location by LID to set LocationID
	var location model.Location
	if err := dbInst.Where("lid = ?", req.Msg.Node.Lid).First(&location).Error; err != nil {
		return nil, err
	}
	node.LocationID = location.ID
	if err := dbInst.Model(&model.Node{}).Where("nid = ?", node.NID).Updates(node).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.UpdateNodeResponse{Success: true}), nil
}

func (h *NodeManagerServiceHandler) DeleteNode(ctx context.Context, req *connect.Request[admin.DeleteNodeRequest]) (*connect.Response[admin.DeleteNodeResponse], error) {
	dbInst := db.Instance()
	if err := dbInst.Where("nid = ?", req.Msg.Nid).Delete(&model.Node{}).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.DeleteNodeResponse{Success: true}), nil
}

func GenerateBackendToken(ctx context.Context, req *connect.Request[admin.GenerateBackendTokenRequest]) (*connect.Response[admin.GenerateBackendTokenResponse], error) {
	dbInst := db.Instance()
	var node model.Node
	if err := dbInst.Where("nid = ?", req.Msg.Nid).First(&node).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("node not found"))
	}

	if node.BackendJTI != nil && *node.BackendJTI != "" && !req.Msg.Regenerate {
		return connect.NewResponse(&admin.GenerateBackendTokenResponse{Success: false}), nil
	}

	backendToken, backendJTI, _, err := session.CreateBackendToken(time.Now())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create backend token"))
	}

	node.BackendJTI = &backendJTI
	if err := dbInst.Save(&node).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to save backend JTI"))
	}

	res := &admin.GenerateBackendTokenResponse{
		Success:      true,
		BackendToken: &backendToken,
	}

	return connect.NewResponse(res), nil
}
