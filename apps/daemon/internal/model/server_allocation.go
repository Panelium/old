package model

type ServerAllocation struct {
	IP       string `gorm:"not null" json:"ip"`
	Port     uint16 `gorm:"not null" json:"port"`
	ServerID uint   `gorm:"index;not null" json:"server_id"`
	Server   Server `json:"server"`
}
