package main

import (
	"github.com/spf13/viper"
	"github.com/vkuzmich/gin-project/internal/app"
	"github.com/vkuzmich/gin-project/internal/http"
	"github.com/vkuzmich/gin-project/pkg/db"
	"log"
)

func main() {
	// Set up Viper to read configuration from .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Get values from Viper
	port := viper.GetString("PORT")
	dbUrl := viper.GetString("DB_URL")

	dbConnection, err := db.Init(dbUrl)
	if err != nil {
		log.Fatalf("Error initializing: %v", err)
	}

	appInstance := app.Build(dbConnection)
	router := http.NewRouter(appInstance)

	err = router.Run(port)
	if err != nil {
		return
	}
}
