package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"panelium/daemon/internal/config"
	"panelium/daemon/internal/model"
	"sync"
)

var (
	initOnce sync.Once
	db       *gorm.DB
)

func Init() error {
	var err error
	initOnce.Do(func() {
		if db, err = gorm.Open(sqlite.Open(config.DatabaseLocation), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}); err != nil {
			return
		}

		if tx := db.Exec("PRAGMA journal_mode=WAL;"); tx.Error != nil {
			err = tx.Error
			return
		}
		if tx := db.Exec("PRAGMA synchronous=NORMAL;"); tx.Error != nil {
			err = tx.Error
			return
		}

		if err = db.AutoMigrate(
			&model.Blueprint{},
			&model.Server{},
			&model.ServerAllocation{},
		); err != nil {
			return
		}
	})

	return err
}

func Instance() *gorm.DB {
	if db == nil {
		panic("database not initialized, call Init() first")
	}
	return db
}
