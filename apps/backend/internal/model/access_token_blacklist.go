package model

import (
	"gorm.io/gorm"
	"time"
)

type AccessTokenBlacklist struct {
	gorm.Model
	JTI        string    `gorm:"uniqueIndex;not null" json:"-"`
	Expiration time.Time `gorm:"not null" json:"expiration"` // current JWT expiration
}
