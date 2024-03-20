package todo_lists

import (
    "github.com/gin-gonic/gin"

    "gorm.io/gorm"
)

// handler is a struct that holds a reference to the database connection.
type handler struct {
    DB *gorm.DB // DB is a reference to the database connection
}

// RegisterRoutes registers routes for the todo list API.
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
    // Create a new handler instance with the provided database connection.
    h := &handler{
        DB: db,
    }

    // Create a new route group for todo tasks under "/books" endpoint.
    routes := r.Group("/todo_tasks")

    // Define routes for CRUD operations on todo tasks.
    routes.POST("/", h.AddTodoTask)       // Create a new todo task
    routes.GET("/", h.GetTodoTasks)       // Get all todo tasks
    routes.GET("/:id", h.GetTodoTask)     // Get a single todo task by ID
    routes.PUT("/:id", h.UpdateTodoTask)  // Update a todo task by ID
    routes.DELETE("/:id", h.DeleteTodoTask) // Delete a todo task by ID
}
