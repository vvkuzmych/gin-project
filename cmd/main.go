package main

import (
	// "fmt"
	// "os"
	// "log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	
	"github.com/vkuzmich/gin-project/pkg/common/db"
	"github.com/vkuzmich/gin-project/pkg/todo_lists"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	port := ":3000"
	dbUrl := "postgres://vladimirkuzmich:mypassword@db:5432/gin_pro1"

	r := gin.Default()
	h := db.Init(dbUrl)

	todo_lists.RegisterRoutes(r, h)
	r.Run(port)
}
