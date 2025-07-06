package admin

import (
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend/admin"
)

func ServerModelToProto(s *model.Server) *admin.Server {
	if s == nil {
		return nil
	}
	uids := make([]string, len(s.Users))
	for i, u := range s.Users {
		uids[i] = u.User.UID
	}
	return &admin.Server{
		Sid:         s.SID,
		Name:        s.Name,
		Description: s.Description,
		OwnerUid:    s.Owner.UID,
		Nid:         s.Node.NID,
		Uids:        uids,
		ResourceLimit: &proto_gen_go.ResourceLimit{
			Cpu:     uint32(s.ResourceLimit.CPU),
			Ram:     uint32(s.ResourceLimit.RAM),
			Swap:    uint32(s.ResourceLimit.SWAP),
			Storage: uint32(s.ResourceLimit.Storage),
		},
		DockerImage: s.DockerImage,
		Bid:         s.BID,
	}
}

func ServerProtoToModel(s *admin.Server) *model.Server {
	if s == nil || s.ResourceLimit == nil {
		return nil
	}
	return &model.Server{
		SID:         s.Sid,
		Name:        s.Name,
		Description: s.Description,
		DockerImage: s.DockerImage,
		BID:         s.Bid,
		ResourceLimit: model.ResourceLimit{
			CPU:     uint(s.ResourceLimit.Cpu),
			RAM:     uint(s.ResourceLimit.Ram),
			SWAP:    uint(s.ResourceLimit.Swap),
			Storage: uint(s.ResourceLimit.Storage),
		},
	}
}
