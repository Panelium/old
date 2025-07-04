package user

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (h *UserServiceHandler) ListUsers(ctx context.Context, req *connect.Request[proto_gen_go.Empty]) (*connect.Response[backend.Users], error) {
	var users []model.User
	if err := db.Instance().Find(&users).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoUsers := make([]*backend.User, len(users))
	for i, user := range users {
		protoUsers[i] = &backend.User{
			Uid: user.UID,
			Data: &backend.UserData{
				Username:  user.Username,
				Email:     user.Email,
				MfaNeeded: user.MFANeeded,
				Admin:     user.Admin,
			},
		}
	}

	response := &backend.Users{
		Users: protoUsers,
	}

	return connect.NewResponse(response), nil
}
