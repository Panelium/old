package model

type ServerAllocation struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	IP     string `gorm:"not null" json:"ip"`
	Port   uint16 `gorm:"not null" json:"port"`
	SID    string `gorm:"index;not null" json:"sid"`
	Server Server `gorm:"foreignKey:SID" json:"server"`
}
