package admin

import (
	"connectrpc.com/connect"
	"context"
	"panelium/proto_gen_go/backend/admin"
)

type UserManagerServiceHandler struct{}

func NewUserManagerServiceHandler() *UserManagerServiceHandler {
	return &UserManagerServiceHandler{}
}

func (h *UserManagerServiceHandler) GetUsers(ctx context.Context, req *connect.Request[admin.GetUsersRequest]) (*connect.Response[admin.GetUsersResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetUsersResponse{}), nil
}

func (h *UserManagerServiceHandler) GetUser(ctx context.Context, req *connect.Request[admin.GetUserRequest]) (*connect.Response[admin.GetUserResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.GetUserResponse{}), nil
}

func (h *UserManagerServiceHandler) CreateUser(ctx context.Context, req *connect.Request[admin.CreateUserRequest]) (*connect.Response[admin.CreateUserResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.CreateUserResponse{}), nil
}

func (h *UserManagerServiceHandler) UpdateUser(ctx context.Context, req *connect.Request[admin.UpdateUserRequest]) (*connect.Response[admin.UpdateUserResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.UpdateUserResponse{}), nil
}

func (h *UserManagerServiceHandler) DeleteUser(ctx context.Context, req *connect.Request[admin.DeleteUserRequest]) (*connect.Response[admin.DeleteUserResponse], error) {
	// TODO: implement
	return connect.NewResponse(&admin.DeleteUserResponse{}), nil
}
