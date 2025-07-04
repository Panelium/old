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

func (h *UserServiceHandler) UpdateUser(ctx context.Context, req *connect.Request[backend.User]) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	var user model.User
	if err := db.Instance().First(&user, req.Msg.Uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	user.Username = req.Msg.Data.Username
	user.Email = req.Msg.Data.Email
	user.MFANeeded = req.Msg.Data.MfaNeeded
	user.Admin = req.Msg.Data.Admin

	if err := db.Instance().Save(&user).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&proto_gen_go.SuccessMessage{Success: true}), nil
}
