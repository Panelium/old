package model

type ServerUser struct {
	ID  uint   `gorm:"primaryKey" json:"id"`
	SID string `gorm:"index,not null;column:sid" json:"sid"`
	UID string `gorm:"index,not null;column:uid" json:"uid"`
}
