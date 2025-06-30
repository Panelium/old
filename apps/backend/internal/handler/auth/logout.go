package auth

import (
	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"context"
	"panelium/backend/internal/middleware"
	"panelium/backend/internal/security/cookies"
	"panelium/backend/internal/security/session"
	"panelium/common/errors"
	proto_gen_go "panelium/proto-gen-go"
)

func (s *AuthServiceHandler) Logout(
	ctx context.Context,
	req *connect.Request[proto_gen_go.LogoutRequest],
) (*connect.Response[proto_gen_go.LogoutResponse], error) {
	sessionInfoData := authn.GetInfo(ctx)
	sessionInfo, ok := sessionInfoData.(*middleware.SessionInfo)
	if !ok || sessionInfo == nil || sessionInfo.SessionID == "" || sessionInfo.UserID == "" {
		return nil, errors.ConnectInvalidCredentials
	}

	err := session.DeleteSession(sessionInfo.SessionID)
	if err != nil {
		res := connect.NewResponse(&proto_gen_go.LogoutResponse{
			Success: false,
		})
		// TODO: log this error?
		return res, connect.NewError(connect.CodeInternal, errors.SessionDeletionFailed)
	}

	res := connect.NewResponse(&proto_gen_go.LogoutResponse{
		Success: true,
	})

	cookies.ClearJWTCookie(res.Header(), "refresh_jwt")
	cookies.ClearJWTCookie(res.Header(), "access_jwt")

	return res, nil
}
