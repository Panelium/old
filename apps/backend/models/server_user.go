package models

import "gorm.io/gorm"

type ServerUser struct {
	gorm.Model
	ServerID uint   `gorm:"index,not null" json:"server_id"`
	Server   Server `json:"server"`
	UserID   uint   `gorm:"index,not null" json:"user_id"`
	User     User   `json:"user"`
}
