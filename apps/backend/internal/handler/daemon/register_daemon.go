package daemon

import (
	"connectrpc.com/connect"
	"context"
	"encoding/base64"
	"errors"
	"panelium/backend/internal/db"
	"panelium/backend/internal/global"
	"panelium/backend/internal/middleware"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *DaemonServiceHandler) RegisterDaemon(
	ctx context.Context,
	req *connect.Request[backend.RegisterDaemonRequest],
) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	daemonInfoData := ctx.Value("panelium_daemon_info")
	daemonInfo, ok := daemonInfoData.(*middleware.DaemonInfo)
	if !ok || daemonInfo == nil || daemonInfo.NID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
	}

	var node *model.Node
	tx := db.Instance().First(&node, "nid = ?", daemonInfo.NID)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("node not found"))
	}

	if node.EncryptedNodeTokenBase64 != nil && *node.EncryptedNodeTokenBase64 != "" {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("node already registered"))
	}

	nodeTokenBytes := []byte(req.Msg.NodeToken)
	encryptedNodeTokenBytes := make([]byte, len(nodeTokenBytes))
	(*global.EncryptionInstance()).Encrypt(encryptedNodeTokenBytes, nodeTokenBytes)
	encryptedNodeTokenBase64 := base64.StdEncoding.EncodeToString(encryptedNodeTokenBytes)

	node.EncryptedNodeTokenBase64 = &encryptedNodeTokenBase64
	if err := db.Instance().Save(node).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to update node token"))
	}

	res := &proto_gen_go.SuccessMessage{
		Success: true,
	}

	return connect.NewResponse(res), nil
}
