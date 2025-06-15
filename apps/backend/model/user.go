package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UID          string       `gorm:"uniqueIndex;not null" json:"uid"`
	Username     string       `gorm:"uniqueIndex;not null" json:"username"`
	Email        string       `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string       `gorm:"not null" json:"-"`
	PasswordSalt string       `gorm:"not null" json:"-"`
	RefreshToken string       `gorm:"not null" json:"-"`
	OwnedServers []Server     `gorm:"foreignKey:OwnerID" json:"owned_servers"`
	Servers      []ServerUser `gorm:"foreignKey:UserID" json:"servers"`
	MFAMethods   []UserMFA    `gorm:"foreignKey:UserID" json:"mfa_methods"`
}
