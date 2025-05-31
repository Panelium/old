package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UID          string       `gorm:"uniqueIndex;not null" json:"uid"`
	Username     string       `gorm:"uniqueIndex;not null" json:"username"`
	Email        string       `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string       `gorm:"not null" json:"-"`
	PasswordSalt string       `gorm:"not null" json:"-"`
	TFAEnabled   bool         `gorm:"default:false" json:"tfa_enabled"`
	TFASecret    string       `gorm:"not null" json:"-"`
	OwnedServers []Server     `gorm:"foreignKey:OwnerID" json:"owned_servers"`
	Servers      []ServerUser `gorm:"foreignKey:UserID" json:"servers"`
}
