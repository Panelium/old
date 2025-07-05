package model

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	SID           string           `gorm:"uniqueIndex;not null" json:"sid"`
	Name          string           `gorm:"not null" json:"name"`
	Description   string           `gorm:"not null" json:"description"`
	OwnerID       uint             `gorm:"index;not null" json:"owner_id"`
	Owner         User             `json:"owner"`
	NodeID        uint             `gorm:"index;not null" json:"node_id"`
	Node          Node             `json:"node"`
	Users         []ServerUser     `gorm:"foreignKey:ServerID" json:"users"`
	Allocations   []NodeAllocation `gorm:"foreignKey:ServerID" json:"allocations"`
	ResourceLimit ResourceLimit    `gorm:"embedded" json:"resource_limit"`
	DockerImage   string           `gorm:"not null" json:"docker_image"`
	BID           string           `gorm:"not null" json:"bid"`
}

type ResourceLimit struct {
	CPU     uint `gorm:"not null" json:"cpu"`     // CPU in percentage (100% = 1 vCore)
	RAM     uint `gorm:"not null" json:"ram"`     // RAM in MB
	SWAP    uint `gorm:"not null" json:"swap"`    // SWAP in MB
	Storage uint `gorm:"not null" json:"storage"` // Storage in MB
}
