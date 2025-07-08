package model

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	LID   string `gorm:"uniqueIndex;not null;column:lid" json:"lid"`
	Name  string `gorm:"uniqueIndex;not null" json:"name"`
	Nodes []Node `gorm:"foreignKey:LocationID" json:"nodes"`
}
