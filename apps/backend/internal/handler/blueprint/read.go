package blueprint

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
)

func (h *BlueprintServiceHandler) ReadBlueprint(ctx context.Context, req *connect.Request[proto_gen_go.SimpleIDMessage]) (*connect.Response[proto_gen_go.SimpleMessage], error) {
	var blueprint model.Blueprint
	if err := db.Instance().Where("bid = ?", req.Msg.Id).First(&blueprint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	blueprintJSON, err := json.Marshal(blueprint)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	response := &proto_gen_go.SimpleMessage{
		Text: string(blueprintJSON),
	}

	return connect.NewResponse(response), nil
}
