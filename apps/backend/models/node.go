package models

import "gorm.io/gorm"

type Node struct {
	gorm.Model
	Name        string       `gorm:"uniqueIndex;not null" json:"name"`
	LocationID  uint         `gorm:"index;not null" json:"location_id"`
	Location    Location     `json:"location"`
	Servers     []Server     `gorm:"foreignKey:NodeID" json:"servers"`
	Allocations []Allocation `gorm:"foreignKey:NodeID" json:"allocations"`
}
