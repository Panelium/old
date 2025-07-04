package blueprint

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
)

func (h *BlueprintServiceHandler) UpdateBlueprint(ctx context.Context, req *connect.Request[proto_gen_go.IDMessage]) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	var blueprint model.Blueprint
	if err := json.Unmarshal([]byte(req.Msg.Text), &blueprint); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	tx := db.Instance().Model(&model.Blueprint{}).Where("bid = ?", req.Msg.Id).Updates(blueprint)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeInternal, tx.Error)
	}

	if err := db.Instance().Save(&blueprint).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&proto_gen_go.SuccessMessage{Success: true}), nil
}
