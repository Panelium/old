package admin

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend/admin"
)

type BlueprintManagerServiceHandler struct{}

func NewBlueprintManagerServiceHandler() *BlueprintManagerServiceHandler {
	return &BlueprintManagerServiceHandler{}
}

func (h *BlueprintManagerServiceHandler) GetBlueprints(ctx context.Context, req *connect.Request[admin.GetBlueprintsRequest]) (*connect.Response[admin.GetBlueprintsResponse], error) {
	dbInst := db.Instance()
	var blueprints []model.Blueprint
	var count int64
	page := req.Msg.Pagination.GetPage()
	pageSize := req.Msg.Pagination.GetPageSize()
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	dbInst.Model(&model.Blueprint{}).Count(&count)
	dbInst.Order("id desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&blueprints)
	resp := &admin.GetBlueprintsResponse{
		Blueprints: make([]*admin.Blueprint, 0, len(blueprints)),
		Pagination: &proto_gen_go.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
	}
	for _, b := range blueprints {
		resp.Blueprints = append(resp.Blueprints, BlueprintModelToProto(&b))
	}
	return connect.NewResponse(resp), nil
}

func (h *BlueprintManagerServiceHandler) GetBlueprint(ctx context.Context, req *connect.Request[admin.GetBlueprintRequest]) (*connect.Response[admin.GetBlueprintResponse], error) {
	dbInst := db.Instance()
	var blueprint model.Blueprint
	if err := dbInst.Where("b_id = ? OR bid = ?", req.Msg.Bid, req.Msg.Bid).First(&blueprint).Error; err != nil {
		return nil, err
	}
	resp := &admin.GetBlueprintResponse{
		Blueprint: BlueprintModelToProto(&blueprint),
	}
	return connect.NewResponse(resp), nil
}

func (h *BlueprintManagerServiceHandler) CreateBlueprint(ctx context.Context, req *connect.Request[admin.CreateBlueprintRequest]) (*connect.Response[admin.CreateBlueprintResponse], error) {
	dbInst := db.Instance()
	var blueprint *model.Blueprint
	if req.Msg.GetBlueprint() != nil {
		blueprint = BlueprintProtoToModel(req.Msg.GetBlueprint())
	} else if req.Msg.GetBlueprintJson() != "" {
		var protoBP admin.Blueprint
		if err := json.Unmarshal([]byte(req.Msg.GetBlueprintJson()), &protoBP); err != nil {
			return nil, err
		}
		blueprint = BlueprintProtoToModel(&protoBP)
	} else {
		return connect.NewResponse(&admin.CreateBlueprintResponse{Success: false}), nil
	}
	if err := dbInst.Create(blueprint).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.CreateBlueprintResponse{Success: true}), nil
}

func (h *BlueprintManagerServiceHandler) UpdateBlueprint(ctx context.Context, req *connect.Request[admin.UpdateBlueprintRequest]) (*connect.Response[admin.UpdateBlueprintResponse], error) {
	dbInst := db.Instance()
	var blueprint *model.Blueprint
	if req.Msg.GetBlueprint() != nil {
		blueprint = BlueprintProtoToModel(req.Msg.GetBlueprint())
	} else if req.Msg.GetBlueprintJson() != "" {
		var protoBP admin.Blueprint
		if err := json.Unmarshal([]byte(req.Msg.GetBlueprintJson()), &protoBP); err != nil {
			return nil, err
		}
		blueprint = BlueprintProtoToModel(&protoBP)
	} else {
		return connect.NewResponse(&admin.UpdateBlueprintResponse{Success: false}), nil
	}
	if err := dbInst.Model(&model.Blueprint{}).Where("bid = ?", blueprint.BID).Updates(blueprint).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.UpdateBlueprintResponse{Success: true}), nil
}

func (h *BlueprintManagerServiceHandler) DeleteBlueprint(ctx context.Context, req *connect.Request[admin.DeleteBlueprintRequest]) (*connect.Response[admin.DeleteBlueprintResponse], error) {
	dbInst := db.Instance()
	if err := dbInst.Where("bid = ?", req.Msg.Bid).Delete(&model.Blueprint{}).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.DeleteBlueprintResponse{Success: true}), nil
}
