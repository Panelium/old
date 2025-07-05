package model

type ServerUser struct {
	SID string `gorm:"index,not null" json:"sid"`
	UID string `gorm:"index,not null" json:"uid"`
}
