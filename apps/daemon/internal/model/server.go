package model

import (
	"panelium/proto_gen_go/daemon"
	"time"
)

type Server struct {
	ID              uint                       `gorm:"primaryKey" json:"id"`
	SID             string                     `gorm:"uniqueIndex;not null" json:"sid"`
	Status          daemon.ServerStatusType    `gorm:"not null" json:"status"`
	TimestampStart  time.Time                  `gorm:"default:null" json:"timestamp_start,omitempty"`
	TimestampEnd    time.Time                  `gorm:"default:null" json:"timestamp_end,omitempty"`
	OfflineReason   daemon.ServerOfflineReason `gorm:"default:null" json:"offline_reason,omitempty"`
	Allocations     []ServerAllocation         `gorm:"foreignKey:ServerID" json:"allocations"`
	ResourceLimit   ResourceLimit              `gorm:"embedded" json:"resource_limit"`
	DockerImage     string                     `gorm:"not null" json:"docker_image"`
	BID             string                     `gorm:"not null" json:"bid"` // Blueprint ID
	Blueprint       Blueprint                  `gorm:"foreignKey:BID" json:"blueprint"`
	ContainerExists bool                       `gorm:"default:false" json:"container_exists"` // Indicates if the server container currently exists in Docker
}

type ResourceLimit struct {
	CPU     uint32 `gorm:"not null" json:"cpu"`     // CPU in percentage (100% = 1 vCore)
	RAM     uint32 `gorm:"not null" json:"ram"`     // RAM in MB
	SWAP    uint32 `gorm:"not null" json:"swap"`    // SWAP in MB
	Storage uint32 `gorm:"not null" json:"storage"` // Storage in MB
}
