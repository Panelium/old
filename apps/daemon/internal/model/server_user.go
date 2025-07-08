package model

type ServerUser struct {
	SID string `gorm:"index,not null;column:sid" json:"sid"`
	UID string `gorm:"index,not null;column:uid" json:"uid"`
}
