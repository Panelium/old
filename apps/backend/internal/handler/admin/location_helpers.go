package admin

import (
	"panelium/backend/internal/model"
	"panelium/proto_gen_go/backend/admin"
)

func LocationModelToProto(l *model.Location) *admin.Location {
	if l == nil {
		return nil
	}
	return &admin.Location{
		Lid:  l.LID,
		Name: l.Name,
	}
}

func LocationProtoToModel(l *admin.Location) *model.Location {
	if l == nil {
		return nil
	}
	return &model.Location{
		LID:  l.Lid,
		Name: l.Name,
	}
}
