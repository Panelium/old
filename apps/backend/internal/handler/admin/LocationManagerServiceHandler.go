package admin

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/common/id"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend/admin"
)

type LocationManagerServiceHandler struct{}

func NewLocationManagerServiceHandler() *LocationManagerServiceHandler {
	return &LocationManagerServiceHandler{}
}

func (h *LocationManagerServiceHandler) GetLocations(ctx context.Context, req *connect.Request[admin.GetLocationsRequest]) (*connect.Response[admin.GetLocationsResponse], error) {
	dbInst := db.Instance()
	var locations []model.Location
	var count int64
	page := req.Msg.Pagination.GetPage()
	pageSize := req.Msg.Pagination.GetPageSize()
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	dbInst.Model(&model.Location{}).Count(&count)
	dbInst.Order("id desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&locations)
	resp := &admin.GetLocationsResponse{
		Locations: make([]*admin.Location, 0, len(locations)),
		Pagination: &proto_gen_go.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    (*uint32)(nil),
		},
	}
	if count > 0 {
		total := uint32(count)
		resp.Pagination.Total = &total
	}
	for _, l := range locations {
		resp.Locations = append(resp.Locations, LocationModelToProto(&l))
	}
	return connect.NewResponse(resp), nil
}

func (h *LocationManagerServiceHandler) GetLocation(ctx context.Context, req *connect.Request[admin.GetLocationRequest]) (*connect.Response[admin.GetLocationResponse], error) {
	dbInst := db.Instance()
	var location model.Location
	if err := dbInst.Where("lid = ?", req.Msg.Lid).First(&location).Error; err != nil {
		return nil, err
	}
	resp := &admin.GetLocationResponse{
		Location: LocationModelToProto(&location),
	}
	return connect.NewResponse(resp), nil
}

func (h *LocationManagerServiceHandler) CreateLocation(ctx context.Context, req *connect.Request[admin.CreateLocationRequest]) (*connect.Response[admin.CreateLocationResponse], error) {
	var err error
	req.Msg.Location.Lid, err = id.New()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create location"))
	}

	dbInst := db.Instance()
	location := LocationProtoToModel(req.Msg.Location)
	if err := dbInst.Create(location).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.CreateLocationResponse{Success: true}), nil
}

func (h *LocationManagerServiceHandler) UpdateLocation(ctx context.Context, req *connect.Request[admin.UpdateLocationRequest]) (*connect.Response[admin.UpdateLocationResponse], error) {
	dbInst := db.Instance()
	location := LocationProtoToModel(req.Msg.Location)
	if err := dbInst.Model(&model.Location{}).Where("lid = ?", location.LID).Updates(location).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.UpdateLocationResponse{Success: true}), nil
}

func (h *LocationManagerServiceHandler) DeleteLocation(ctx context.Context, req *connect.Request[admin.DeleteLocationRequest]) (*connect.Response[admin.DeleteLocationResponse], error) {
	dbInst := db.Instance()
	if err := dbInst.Where("lid = ?", req.Msg.Lid).Delete(&model.Location{}).Error; err != nil {
		return nil, err
	}
	return connect.NewResponse(&admin.DeleteLocationResponse{Success: true}), nil
}
