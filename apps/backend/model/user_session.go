package model

type UserSession struct {
	SessionID    string `gorm:"uniqueIndex;not null" json:"session_id"`
	UserID       string `gorm:"index;not null" json:"user_id"`
	User         User   `json:"user"`
	RefreshToken string `gorm:"not null" json:"-"`
}
