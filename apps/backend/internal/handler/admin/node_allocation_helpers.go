package admin

import (
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend/admin"
)

func NodeAllocationModelToProto(na *model.NodeAllocation) *admin.NodeAllocation {
	if na == nil {
		return nil
	}
	return &admin.NodeAllocation{
		Id:  uint32(na.ID),
		Nid: na.Node.NID,
		Sid: func() *string {
			if na.ServerID != nil && *na.ServerID != 0 {
				s := na.Server.SID
				return &s
			}
			return nil
		}(),
		IpAllocation: &proto_gen_go.IPAllocation{
			Ip:   na.IP,
			Port: uint32(na.Port),
		},
	}
}

func NodeAllocationProtoToModel(na *admin.NodeAllocation) *model.NodeAllocation {
	if na == nil || na.IpAllocation == nil {
		return nil
	}

	dbNa := &model.NodeAllocation{
		IP:   na.IpAllocation.Ip,
		Port: uint16(na.IpAllocation.Port),
	}

	if na.Sid != nil && *na.Sid != "" {
		var server model.Server
		tx := db.Instance().Where("sid = ?", *na.Sid).First(&server)
		if tx.Error != nil || tx.RowsAffected == 0 {
			return nil
		}

		dbNa.ServerID = &server.ID
	}

	var node model.Node
	if err := db.Instance().Where("nid = ?", na.Nid).First(&node).Error; err != nil {
		return nil
	}

	dbNa.NodeID = node.ID

	return dbNa
}
