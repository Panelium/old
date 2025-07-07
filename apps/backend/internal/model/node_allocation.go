package model

import "gorm.io/gorm"

type NodeAllocation struct {
	gorm.Model
	NodeID   uint   `gorm:"index,not null" json:"node_id"`
	Node     Node   `json:"node"`
	IP       string `gorm:"not null" json:"ip"`
	Port     uint16 `gorm:"not null" json:"port"`
	ServerID *uint  `gorm:"index" json:"server_id"` // Nullable, does not have to be assigned to a server
	Server   Server `json:"server,omitempty"`       // Use omitempty to avoid null in JSON if ServerID is not set
}
