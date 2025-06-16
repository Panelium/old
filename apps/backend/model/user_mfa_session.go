package model

type UserMFASession struct {
	SessionID string `gorm:"uniqueIndex;not null" json:"session_id"`
	UserID    string `gorm:"index;not null" json:"user_id"`
	User      User   `json:"user"`
	Token     string `gorm:"not null" json:"-"`
}
