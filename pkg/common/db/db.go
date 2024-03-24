package db

import (
	// "fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/vkuzmich/gin-project/pkg/common/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	
    if err != nil {
        log.Fatalln(err)
    }

    db.AutoMigrate(&models.TodoTask{})

    return db
}
