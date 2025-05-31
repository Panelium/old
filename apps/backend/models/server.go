package models

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	SID         string       `gorm:"uniqueIndex;not null" json:"sid"`
	Name        string       `gorm:"not null" json:"name"`
	OwnerID     uint         `gorm:"index;not null" json:"owner_id"`
	Owner       User         `json:"owner"`
	NodeID      uint         `gorm:"index;not null" json:"node_id"`
	Node        Node         `json:"node"`
	Users       []ServerUser `gorm:"foreignKey:ServerID" json:"users"`
	Allocations []Allocation `gorm:"foreignKey:ServerID" json:"allocations"`
}
