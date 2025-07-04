package model

type ServerUser struct {
	SID    string `gorm:"index,not null" json:"sid"`
	UserID string `gorm:"index,not null" json:"user_id"`
}
