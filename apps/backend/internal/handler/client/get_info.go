package client

import (
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/db"
	"panelium/backend/internal/middleware"
	"panelium/backend/internal/model"
	"panelium/common/errors"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *ClientServiceHandler) GetInfo(ctx context.Context, req *connect.Request[proto_gen_go.Empty]) (*connect.Response[backend.ClientInfo], error) {
	sessionInfoData := ctx.Value("panelium_session_info")
	sessionInfo, ok := sessionInfoData.(*middleware.SessionInfo)
	if !ok || sessionInfo == nil || sessionInfo.SessionID == "" || sessionInfo.UserID == "" {
		return nil, errors.ConnectInvalidCredentials
	}

	var user *model.User
	tx := db.Instance().First(&user, "uid = ?", sessionInfo.UserID)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.UserNotFound)
	}

	res := &backend.ClientInfo{
		Uid:      sessionInfo.UserID,
		Username: user.Username,
		Email:    user.Email,
		Admin:    user.Admin,
	}

	return connect.NewResponse(res), nil
}
