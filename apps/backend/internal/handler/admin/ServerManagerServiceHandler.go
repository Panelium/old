package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend/admin"
)

type ServerManagerServiceHandler struct{}

func NewServerManagerServiceHandler() *ServerManagerServiceHandler {
	return &ServerManagerServiceHandler{}
}

func (h *ServerManagerServiceHandler) GetServers(ctx context.Context, req *connect.Request[admin.GetServersRequest]) (*connect.Response[admin.GetServersResponse], error) {
	dbInst := db.Instance()
	var servers []model.Server
	var count int64
	page := req.Msg.Pagination.GetPage()
	pageSize := req.Msg.Pagination.GetPageSize()
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	query := dbInst.Model(&model.Server{}).Preload("Owner").Preload("Node").Preload("Users.User")
	if req.Msg.Nid != nil && *req.Msg.Nid != "" {
		query = query.Joins("Node").Where("nodes.nid = ?", *req.Msg.Nid)
	}
	if req.Msg.OwnerUid != nil && *req.Msg.OwnerUid != "" {
		query = query.Joins("Owner").Where("users.uid = ?", *req.Msg.OwnerUid)
	}
	if req.Msg.Bid != nil && *req.Msg.Bid != "" {
		query = query.Where("bid = ?", *req.Msg.Bid)
	}
	if req.Msg.AccessUid != nil && *req.Msg.AccessUid != "" {
		query = query.Joins("Users.User").Where("users.uid = ?", *req.Msg.AccessUid)
	}
	if req.Msg.Uid != nil && *req.Msg.Uid != "" {
		query = query.Joins("Owner").Joins("Users.User").Where("users.uid = ? OR owner_id = (SELECT id FROM users WHERE uid = ?)", *req.Msg.Uid, *req.Msg.Uid)
	}
	query.Count(&count)
	query.Order("id desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&servers)
	resp := &admin.GetServersResponse{
		Servers: make([]*admin.Server, 0, len(servers)),
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
	for _, s := range servers {
		resp.Servers = append(resp.Servers, ServerModelToProto(&s))
	}
	return connect.NewResponse(resp), nil
}

func (h *ServerManagerServiceHandler) GetServer(ctx context.Context, req *connect.Request[admin.GetServerRequest]) (*connect.Response[admin.GetServerResponse], error) {
	dbInst := db.Instance()
	var server model.Server
	if err := dbInst.Preload("Owner").Preload("Node").Preload("Users.User").Where("sid = ?", req.Msg.Sid).First(&server).Error; err != nil {
		return nil, err
	}
	resp := &admin.GetServerResponse{
		Server: ServerModelToProto(&server),
	}
	return connect.NewResponse(resp), nil
}

func (h *ServerManagerServiceHandler) CreateServer(ctx context.Context, req *connect.Request[admin.CreateServerRequest]) (*connect.Response[admin.CreateServerResponse], error) {
	if req.Msg.Server.Sid != "" {
		req.Msg.Server.Sid = ""
	}
	dbInst := db.Instance()
	server := ServerProtoToModel(req.Msg.Server)
	// Find owner by UID
	var owner model.User
	if err := dbInst.Where("uid = ?", req.Msg.Server.OwnerUid).First(&owner).Error; err != nil {
		return nil, err
	}
	server.OwnerID = owner.ID
	// Find node by NID
	var node model.Node
	if err := dbInst.Where("nid = ?", req.Msg.Server.Nid).First(&node).Error; err != nil {
		return nil, err
	}
	server.NodeID = node.ID
	if err := dbInst.Create(server).Error; err != nil {
		return nil, err
	}
	// Add users if provided
	for _, uid := range req.Msg.Server.Uids {
		var user model.User
		if err := dbInst.Where("uid = ?", uid).First(&user).Error; err == nil {
			_ = dbInst.Create(&model.ServerUser{ServerID: server.ID, UserID: user.ID}).Error
		}
	}
	return connect.NewResponse(&admin.CreateServerResponse{Success: true}), nil
}

func (h *ServerManagerServiceHandler) UpdateServer(ctx context.Context, req *connect.Request[admin.UpdateServerRequest]) (*connect.Response[admin.UpdateServerResponse], error) {
	dbInst := db.Instance()
	server := ServerProtoToModel(req.Msg.Server)
	// Find owner by UID
	var owner model.User
	if err := dbInst.Where("uid = ?", req.Msg.Server.OwnerUid).First(&owner).Error; err != nil {
		return nil, err
	}
	server.OwnerID = owner.ID
	// Find node by NID
	var node model.Node
	if err := dbInst.Where("nid = ?", req.Msg.Server.Nid).First(&node).Error; err != nil {
		return nil, err
	}
	server.NodeID = node.ID
	if err := dbInst.Model(&model.Server{}).Where("sid = ?", server.SID).Updates(server).Error; err != nil {
		return nil, err
	}
	// Update users
	dbInst.Where("server_id = ?", server.ID).Delete(&model.ServerUser{})
	for _, uid := range req.Msg.Server.Uids {
		var user model.User
		if err := dbInst.Where("uid = ?", uid).First(&user).Error; err == nil {
			_ = dbInst.Create(&model.ServerUser{ServerID: server.ID, UserID: user.ID}).Error
		}
	}
	return connect.NewResponse(&admin.UpdateServerResponse{Success: true}), nil
}

func (h *ServerManagerServiceHandler) DeleteServer(ctx context.Context, req *connect.Request[admin.DeleteServerRequest]) (*connect.Response[admin.DeleteServerResponse], error) {
	dbInst := db.Instance()
	if err := dbInst.Where("sid = ?", req.Msg.Sid).Delete(&model.Server{}).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.DeleteServerResponse{Success: true}), nil
}
