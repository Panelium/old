package user

import (
	"connectrpc.com/connect"
	"context"
	"gorm.io/gorm"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (h *UserServiceHandler) ReadUser(ctx context.Context, req *connect.Request[proto_gen_go.SimpleIDMessage]) (*connect.Response[backend.User], error) {
	var user model.User
	if err := db.Instance().Where("uid = ?", req.Msg.Id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	response := &backend.User{
		Uid: user.UID,
		Data: &backend.UserData{
			Username:  user.Username,
			Email:     user.Email,
			MfaNeeded: user.MFANeeded,
			Admin:     user.Admin,
		},
	}

	return connect.NewResponse(response), nil
}

func (h *UserServiceHandler) ReadUserByUsername(ctx context.Context, req *connect.Request[proto_gen_go.SimpleMessage]) (*connect.Response[backend.User], error) {
	var user model.User
	if err := db.Instance().Where("username = ?", req.Msg.Text).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	response := &backend.User{
		Uid: user.UID,
		Data: &backend.UserData{
			Username:  user.Username,
			Email:     user.Email,
			MfaNeeded: user.MFANeeded,
			Admin:     user.Admin,
		},
	}

	return connect.NewResponse(response), nil
}

func (h *UserServiceHandler) ReadUserByEmail(ctx context.Context, req *connect.Request[proto_gen_go.SimpleMessage]) (*connect.Response[backend.User], error) {
	var user model.User
	if err := db.Instance().Where("email = ?", req.Msg.Text).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	response := &backend.User{
		Uid: user.UID,
		Data: &backend.UserData{
			Username:  user.Username,
			Email:     user.Email,
			MfaNeeded: user.MFANeeded,
			Admin:     user.Admin,
		},
	}

	return connect.NewResponse(response), nil
}
