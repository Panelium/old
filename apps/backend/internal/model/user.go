package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UID          string           `gorm:"uniqueIndex;not null;column:uid" json:"uid"`
	Username     string           `gorm:"uniqueIndex;not null" json:"username"`
	Email        string           `gorm:"uniqueIndex;not null" json:"email"`
	Admin        bool             `gorm:"not null;default:false" json:"admin"`
	PasswordHash string           `gorm:"not null" json:"-"`
	PasswordSalt string           `gorm:"not null" json:"-"`
	MFANeeded    bool             `gorm:"not null" json:"mfa_needed"`
	MFAMethods   []UserMFA        `gorm:"foreignKey:UserID" json:"mfa_methods"`
	MFASessions  []UserMFASession `gorm:"foreignKey:UserID" json:"mfa_sessions"`
	Sessions     []UserSession    `gorm:"foreignKey:UserID" json:"sessions"`
	OwnedServers []Server         `gorm:"foreignKey:OwnerID" json:"owned_servers"`
	Servers      []ServerUser     `gorm:"foreignKey:UserID" json:"servers"`
}
