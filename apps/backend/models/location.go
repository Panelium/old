package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex;not null" json:"name"`
	Nodes []Node `gorm:"foreignKey:LocationID" json:"nodes"`
}
