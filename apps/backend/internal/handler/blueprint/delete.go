package blueprint

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
)

func (h *BlueprintServiceHandler) DeleteBlueprint(ctx context.Context, req *connect.Request[proto_gen_go.SimpleIDMessage]) (*connect.Response[proto_gen_go.SuccessMessage], error) {
	tx := db.Instance().Find(&model.Server{}, "bid = ?", req.Msg.Id)
	if tx.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, tx.Error)
	}
	if tx.RowsAffected > 0 {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("blueprint is in use by servers"))
	}

	if err := db.Instance().Where("bid = ?", req.Msg.Id).Delete(&model.Blueprint{}).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&proto_gen_go.SuccessMessage{Success: true}), nil
}
