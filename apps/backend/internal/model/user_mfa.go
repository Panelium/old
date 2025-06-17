package model

import "gorm.io/gorm"

type UserMFA struct {
	gorm.Model
	UserID  uint   `gorm:"index,not null" json:"user_id"`
	User    User   `json:"user"`
	Type    string `gorm:"not null" json:"type"`    // e.g., "totp", "sms", "email"
	Value   string `json:"value"`                   // e.g., phone number, email address
	Secret  string `json:"-"`                       // Secret used for TOTP or other methods, not exposed in JSON
	Enabled bool   `gorm:"not null" json:"enabled"` // Whether the MFA method is enabled
}
