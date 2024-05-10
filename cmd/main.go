package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/vkuzmich/gin-project/db"
	"github.com/vkuzmich/gin-project/middleware"
	"github.com/vkuzmich/gin-project/routes"
	"log"
)

func main() {
	// Set up Viper to read configuration from .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	//
	// Get values from Viper
	port := viper.GetString("PORT")
	dbUrl := viper.GetString("DB_URL")

	r := gin.Default()
	r.Use(middleware.RouteMiddleware())
	h, err := db.Init(dbUrl)
	if err != nil {
		log.Fatalf("Error initializing: %v", err)
	}

	routes.RegisterRoutes(r, h)
	r.Run(port)
}
