package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend/admin"
)

type UserManagerServiceHandler struct{}

func NewUserManagerServiceHandler() *UserManagerServiceHandler {
	return &UserManagerServiceHandler{}
}

func (h *UserManagerServiceHandler) GetUsers(ctx context.Context, req *connect.Request[admin.GetUsersRequest]) (*connect.Response[admin.GetUsersResponse], error) {
	dbInst := db.Instance()
	var users []model.User
	var count int64
	page := req.Msg.Pagination.GetPage()
	pageSize := req.Msg.Pagination.GetPageSize()
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	dbInst.Model(&model.User{}).Count(&count)
	dbInst.Order("id desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&users)
	resp := &admin.GetUsersResponse{
		Users: make([]*admin.User, 0, len(users)),
		Pagination: &proto_gen_go.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
	}
	for _, u := range users {
		resp.Users = append(resp.Users, UserModelToProto(&u))
	}
	return connect.NewResponse(resp), nil
}

func (h *UserManagerServiceHandler) GetUser(ctx context.Context, req *connect.Request[admin.GetUserRequest]) (*connect.Response[admin.GetUserResponse], error) {
	dbInst := db.Instance()
	var user model.User
	if err := dbInst.Where("uid = ?", req.Msg.Uid).First(&user).Error; err != nil {
		return nil, err
	}
	resp := &admin.GetUserResponse{
		User: UserModelToProto(&user),
	}
	return connect.NewResponse(resp), nil
}

func (h *UserManagerServiceHandler) CreateUser(ctx context.Context, req *connect.Request[admin.CreateUserRequest]) (*connect.Response[admin.CreateUserResponse], error) {
	if req.Msg.User.Uid != "" {
		req.Msg.User.Uid = ""
	}
	dbInst := db.Instance()
	user := UserProtoToModel(req.Msg.User)
	if err := dbInst.Create(user).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.CreateUserResponse{Success: true}), nil
}

func (h *UserManagerServiceHandler) UpdateUser(ctx context.Context, req *connect.Request[admin.UpdateUserRequest]) (*connect.Response[admin.UpdateUserResponse], error) {
	dbInst := db.Instance()
	user := UserProtoToModel(req.Msg.User)
	if err := dbInst.Model(&model.User{}).Where("uid = ?", user.UID).Updates(user).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.UpdateUserResponse{Success: true}), nil
}

func (h *UserManagerServiceHandler) DeleteUser(ctx context.Context, req *connect.Request[admin.DeleteUserRequest]) (*connect.Response[admin.DeleteUserResponse], error) {
	dbInst := db.Instance()
	if err := dbInst.Where("uid = ?", req.Msg.Uid).Delete(&model.User{}).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.DeleteUserResponse{Success: true}), nil
}
