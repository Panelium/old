package model

import "time"

type UserSession struct {
	SessionID  string    `gorm:"uniqueIndex;not null" json:"session_id"`
	UserID     string    `gorm:"index;not null" json:"user_id"`
	User       User      `json:"user"`
	AccessJTI  string    `gorm:"uniqueIndex;not null" json:"-"`
	RefreshJTI string    `gorm:"uniqueIndex;not null" json:"-"`
	Expiration time.Time `gorm:"not null" json:"expiration"` // same as current refresh JWT expiration
}
