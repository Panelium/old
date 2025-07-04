package model

type ServerUser struct {
	ServerID uint   `gorm:"index,not null" json:"server_id"`
	UserID   string `gorm:"index,not null" json:"user_id"`
}
