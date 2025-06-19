package global

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"panelium/backend/internal/model"
)

var DB *gorm.DB

const Pepper string = "s0m3R4nd0mP3pp3r"     // TODO!: move this to a config
const JWTSecret string = "0th3rR4nd0mS3cr3t" // TODO!: move this to a config

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
		&model.AccessTokenBlacklist{},
		&model.Blueprint{},
		&model.Location{},
		&model.Node{},
		&model.NodeAllocation{},
		&model.Server{},
		&model.ServerUser{},
		&model.User{},
		&model.UserMFA{},
		&model.UserMFASession{},
		&model.UserSession{},
	)
	if err != nil {
		return err
	}
	return nil
}
