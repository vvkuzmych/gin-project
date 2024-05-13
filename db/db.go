package db

import (
	"errors"
	"fmt"
	"github.com/vkuzmich/gin-project/internal/model"
	"strings"

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
	return db.AutoMigrate(&model.TodoTask{})
}

func ConnectionToDB(url string) (*gorm.DB, error) {
	// Attempt to open the database connection
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		// Check if the error is related to database does not exist
		if strings.Contains(err.Error(), `database "gin_pron" does not exist`) {
			// Handle the error accordingly without printing it
			return nil, fmt.Errorf("failed to initialize database: database does not exist")
		}
		// Return the original error if it's not related to database does not exist
		return nil, err
	}

	// Return the database connection and nil error if successful
	return db, nil
}
