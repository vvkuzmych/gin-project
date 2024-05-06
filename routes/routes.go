package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vkuzmich/gin-project/handler"
	"gorm.io/gorm"
)

// handler is a struct that holds a reference to the database connection.
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &handler.MyDatabase{DB: db}

	routes := r.Group("/todo_tasks")

	routes.POST("/", h.AddTodoTask)
	routes.GET("/", h.GetTodoTasks)
	routes.GET("/:id", h.GetTodoTask)
	routes.PUT("/:id", h.UpdateTodoTask)
	routes.DELETE("/:id", h.DeleteTodoTask)
}
