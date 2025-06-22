package model

import "time"

type UserMFASession struct {
	SessionID  string    `gorm:"uniqueIndex;not null" json:"session_id"`
	UserID     string    `gorm:"index;not null" json:"user_id"`
	User       User      `json:"user"`
	JTI        string    `gorm:"uniqueIndex;not null" json:"-"`
	Expiration time.Time `gorm:"not null" json:"expiration"` // current JWT expiration
}
