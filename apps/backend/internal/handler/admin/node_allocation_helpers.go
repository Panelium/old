package admin

import (
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
			if na.ServerID != 0 {
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
	return &model.NodeAllocation{
		IP:   na.IpAllocation.Ip,
		Port: uint16(na.IpAllocation.Port),
	}
}
