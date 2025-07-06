package admin

import (
	"panelium/backend/internal/model"
	"panelium/proto_gen_go/backend/admin"
)

func NodeModelToProto(n *model.Node) *admin.Node {
	if n == nil {
		return nil
	}
	return &admin.Node{
		Nid:        n.NID,
		Name:       n.Name,
		Fqdn:       n.FQDN,
		DaemonPort: uint32(n.DaemonPort),
		Https:      n.HTTPS,
		Lid:        n.Location.LID,
		MaxCpu:     uint32(n.MaxCPU),
		MaxRam:     uint32(n.MaxRAM),
		MaxSwap:    uint32(n.MaxSWAP),
		MaxStorage: uint32(n.MaxStorage),
	}
}

func NodeProtoToModel(n *admin.Node) *model.Node {
	if n == nil {
		return nil
	}
	return &model.Node{
		NID:        n.Nid,
		Name:       n.Name,
		FQDN:       n.Fqdn,
		DaemonPort: uint(n.DaemonPort),
		HTTPS:      n.Https,
		MaxCPU:     uint(n.MaxCpu),
		MaxRAM:     uint(n.MaxRam),
		MaxSWAP:    uint(n.MaxSwap),
		MaxStorage: uint(n.MaxStorage),
	}
}
