package global

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	model2 "panelium/backend/internal/model"
)

var DB *gorm.DB
var Pepper string = "s0m3R4nd0mP3pp3r"     // TODO!: move this to a config
var JWTSecret string = "0th3rR4nd0mS3cr3t" // TODO!: move this to a config

func Init() error {
	var err error
	// TODO: sqlite is temporary, change in the future
	DB, err = gorm.Open(sqlite.Open("file:panelium.db?_journal_mode=WAL"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(
		&model2.AccessTokenBlacklist{},
		&model2.Blueprint{},
		&model2.Location{},
		&model2.Node{},
		&model2.NodeAllocation{},
		&model2.Server{},
		&model2.ServerUser{},
		&model2.User{},
		&model2.UserMFA{},
	)
	if err != nil {
		return err
	}
	return nil
}
