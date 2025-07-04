package blueprint

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
)

func (h *BlueprintServiceHandler) ListBlueprints(ctx context.Context, req *connect.Request[proto_gen_go.Empty]) (*connect.Response[proto_gen_go.SimpleMessage], error) {
	var blueprints []model.Blueprint
	if err := db.Instance().Find(&blueprints).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	blueprintJSON, err := json.Marshal(blueprints)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	response := &proto_gen_go.SimpleMessage{
		Text: string(blueprintJSON),
	}

	return connect.NewResponse(response), nil
}
