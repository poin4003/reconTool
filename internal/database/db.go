package database

import (
	"log"
	"time"

	"reconTool/global"
	"reconTool/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() *gorm.DB {
	dsn := "recon.db?_busy_timeout=5000&_journal_mode=WAL"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Init database initialization error: %v", err)
	}

	log.Println("Initializing Successfully")

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Init database get DB error: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	global.DB = db
	return db
}

func DBMigrator(db *gorm.DB) error {
	log.Println("Running AutoMigrate...")

	err := db.AutoMigrate(
		&model.ConfigModel{},
		&model.SimModel{},
	)

	if err != nil {
		return err
	}

	return nil
}
