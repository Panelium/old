package user

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/common/id"
	"panelium/proto_gen_go/backend"
)

func (h *UserServiceHandler) CreateUser(ctx context.Context, req *connect.Request[backend.UserData]) (*connect.Response[backend.User], error) {
	uid, err := id.New()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	user := &model.User{
		UID:       uid,
		Username:  req.Msg.Username,
		Email:     req.Msg.Email,
		MFANeeded: req.Msg.MfaNeeded,
		Admin:     req.Msg.Admin,
	}

	if err := db.Instance().Create(user).Error; err != nil {
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
