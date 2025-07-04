package user

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
)

func (h *UserServiceHandler) DeleteUser(ctx context.Context, req *connect.Request[proto_gen_go.SimpleIDMessage]) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	// TODO: cascade delete user's servers, sessions, mfa, ...

	if err := db.Instance().Where("uid = ?", req.Msg.Id).Delete(&model.User{}).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&proto_gen_go.SuccessMessage{Success: true}), nil
}
