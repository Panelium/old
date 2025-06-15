package model

import (
	"gorm.io/gorm"
	"time"
)

type AccessTokenBlacklist struct {
	gorm.Model
	TokenHash  string    `gorm:"uniqueIndex;not null" json:"token_hash"`
	Expiration time.Time `gorm:"not null" json:"expiration"`
}
