package model

import proto_gen_go "panelium/proto-gen-go"

type Server struct {
	ID            uint                          `gorm:"primaryKey" json:"id"`
	SID           string                        `gorm:"uniqueIndex;not null" json:"sid"`
	Status        proto_gen_go.ServerStatusType `gorm:"not null" json:"status"`
	Allocations   []ServerAllocation            `gorm:"foreignKey:ServerID" json:"allocations"`
	ResourceLimit ResourceLimit                 `gorm:"embedded" json:"resource_limit"`
	DockerImage   string                        `gorm:"not null" json:"docker_image"`
	BID           string                        `gorm:"not null" json:"bid"` // Blueprint ID
	Blueprint     Blueprint                     `gorm:"foreignKey:BID" json:"blueprint"`
}

type ResourceLimit struct {
	CPU     uint `gorm:"not null" json:"cpu"`     // CPU in percentage (100% = 1 vCore)
	RAM     uint `gorm:"not null" json:"ram"`     // RAM in MB
	SWAP    uint `gorm:"not null" json:"swap"`    // SWAP in MB
	Storage uint `gorm:"not null" json:"storage"` // Storage in MB
}
