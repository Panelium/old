package model

import "gorm.io/gorm"

type Node struct {
	gorm.Model
	Name        string           `gorm:"uniqueIndex;not null" json:"name"`
	FQDN        string           `gorm:"uniqueIndex;not null" json:"fqdn"`
	LocationID  uint             `gorm:"index;not null" json:"location_id"`
	Location    Location         `json:"location"`
	Servers     []Server         `gorm:"foreignKey:NodeID" json:"servers"`
	Allocations []NodeAllocation `gorm:"foreignKey:NodeID" json:"allocations"`
	MaxRAM      uint             `gorm:"not null" json:"max_ram"`     // Maximum RAM in MB
	MaxStorage  uint             `gorm:"not null" json:"max_storage"` // Maximum Storage in MB
	NodeToken   string           `gorm:"not null" json:"node_token"`  // Encrypted node token (backend->daemon communication)
	BackendJTI  string           `gorm:"not null" json:"backend_jti"` // JWT ID of the backend token (daemon->backend communication)
}
