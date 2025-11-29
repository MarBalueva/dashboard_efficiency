package db

import (
	"log"

	"github.com/MarBalueva/dashboard_efficiency/internal/config"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := cfg.PostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func Migrate() error {
	if DB == nil {
		return nil
	}

	err := DB.AutoMigrate(
		&models.Employee{},
		&models.User{},
	)

	if err != nil {
		log.Println("migration error:", err)
		return err
	}

	return nil
}
