package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	logger := logger.Default.LogMode((logger.Silent))
	connection, err := gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	DB = connection
}

func AutoMigrate(models ...interface{}) error {
	return DB.AutoMigrate(models...)
}
