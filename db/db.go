package db

import (
	"errors"
	"github.com/vkuzmich/gin-project/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ConnectingToDB = ConnectionToDB
	AutoMigrations = AutoMigration
)

func Init(url string) (*gorm.DB, error) {
	db, err := ConnectingToDB(url)
	if err != nil {
		return nil, err
	}

	err = AutoMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigration(db *gorm.DB) error {
	if db == nil {
		return errors.New("nil database connection")
	}
	return db.AutoMigrate(&models.TodoTask{})
}

func ConnectionToDB(url string) (*gorm.DB, error) {
	// Attempt to open the database connection
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		// Return the error if opening the connection fails
		return nil, err
	}

	// Return the database connection and nil error if successful
	return db, nil
}
