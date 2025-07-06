package model

import "gorm.io/gorm"

type Node struct {
	gorm.Model
	NID         string           `gorm:"uniqueIndex;not null" json:"nid"` // Node ID
	Name        string           `gorm:"uniqueIndex;not null" json:"name"`
	FQDN        string           `gorm:"uniqueIndex;not null" json:"fqdn"`
	DaemonPort  uint             `gorm:"not null" json:"daemon_port"`
	HTTPS       bool             `gorm:"not null" json:"https"` // Whether the node uses HTTPS for the daemon communication
	LocationID  uint             `gorm:"index;not null" json:"location_id"`
	Location    Location         `json:"location"`
	Servers     []Server         `gorm:"foreignKey:NodeID" json:"servers"`
	Allocations []NodeAllocation `gorm:"foreignKey:NodeID" json:"allocations"`
	MaxCPU      uint             `gorm:"not null" json:"max_cpu"`     // Maximum CPU in % (100% = 1 vCore)
	MaxRAM      uint             `gorm:"not null" json:"max_ram"`     // Maximum RAM in MB
	MaxSWAP     uint             `gorm:"not null" json:"max_swap"`    // Maximum Swap in MB
	MaxStorage  uint             `gorm:"not null" json:"max_storage"` // Maximum Storage in MB
	NodeToken   *string          `json:"node_token"`                  // Encrypted node token (backend->daemon communication)
	BackendJTI  *string          `json:"backend_jti"`                 // JWT ID of the backend token (daemon->backend communication)
}
