package blueprint

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
)

func (h *BlueprintServiceHandler) CreateBlueprint(ctx context.Context, req *connect.Request[proto_gen_go.SimpleMessage]) (*connect.Response[proto_gen_go.SimpleIDMessage], error) {
	blueprint := &model.Blueprint{}
	if err := json.Unmarshal([]byte(req.Msg.Text), &blueprint); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := db.Instance().Create(blueprint).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	response := &proto_gen_go.SimpleIDMessage{
		Id: blueprint.BID,
	}

	return connect.NewResponse(response), nil
}
