package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vkuzmich/gin-project/internal/app"
	"github.com/vkuzmich/gin-project/middleware"
	"github.com/vkuzmich/gin-project/routes"
)

func NewRouter(a app.Interface) *gin.Engine {

	var (
		todoTaskService = a.TodoTaskService()
	)
	router := gin.Default()
	router.Use(middleware.RouteMiddleware())

	v := router.Group("")
	fmt.Println("Starting application...v", v)

	routes.RegisterTodoTaskHandlers(v, todoTaskService)
	return router
}
