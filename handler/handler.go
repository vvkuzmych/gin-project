package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vkuzmich/gin-project/internal/models"
	"gorm.io/gorm"
	"net/http"
)

// AddTodoTaskRequestBody represents the structure of the request body
// expected when adding a new todo task.
type TodoTaskRequestBody struct {
	Title       string `json:"title"`       // Title of the todo task
	Description string `json:"description"` // Description of the todo task
	State       bool   `json:"state"`       // State of the todo task (completed or not)
}
type TodoHandler interface {
	AddTodoTask(c *gin.Context)
	DeleteTodoTask(c *gin.Context)
	GetTodoTask(c *gin.Context)
	GetTodoTasks(c *gin.Context)
	UpdateTodoTask(c *gin.Context)
}

// MyDatabase represents a database connection.
type MyDatabase struct {
	DB *gorm.DB // DB is a reference to the database connection
}

// AddTodoTask is a handler function for adding a new todo task.
func (h *MyDatabase) AddTodoTask(c *gin.Context) {
	body := TodoTaskRequestBody{}

	// Receive request body
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var todoTask models.TodoTask

	// Assign values from request body to todoTask
	todoTask.Title = body.Title
	todoTask.Description = body.Description
	todoTask.State = body.State

	// Create todo task in the database
	if result := h.DB.Create(&todoTask); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	// Respond with the created todo task
	c.JSON(http.StatusCreated, &todoTask)
}

// DeleteTodoTask is a handler function for deleting a todo task.
func (h *MyDatabase) DeleteTodoTask(c *gin.Context) {
	// Extract the ID parameter from the request URL.
	id := c.Param("id")

	// Declare a variable to store the retrieved todo task.
	var todoTask models.TodoTask

	// Retrieve the todo task from the database by its ID.
	if result := h.DB.First(&todoTask, id); result.Error != nil {
		// Abort the request with an error if retrieval fails.
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	// Delete the retrieved todo task from the database.
	h.DB.Delete(&todoTask)

	// Respond with a success status.
	c.Status(http.StatusOK)
}

// GetTodoTask is a handler function for retrieving a single todo task by ID.
func (h *MyDatabase) GetTodoTask(c *gin.Context) {
	// Extract the ID parameter from the request URL.
	id := c.Param("id")

	// Declare a variable to store the retrieved todo task.
	var todoTask models.TodoTask

	// Retrieve the todo task from the database by its ID.
	if result := h.DB.First(&todoTask, id); result.Error != nil {
		// Abort the request with an error if retrieval fails.
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	// Respond with the retrieved todo task.
	c.JSON(http.StatusOK, &todoTask)
}

// GetTodoTasks is a handler function for retrieving todo tasks.
func (h *MyDatabase) GetTodoTasks(c *gin.Context) {
	// Define a slice to store todo tasks.
	var todoTasks []models.TodoTask

	// Retrieve todo tasks from the database.
	if result := h.DB.Find(&todoTasks); result.Error != nil {
		// Abort the request with an error if retrieval fails.
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	// Respond with the retrieved todo tasks.
	c.JSON(http.StatusOK, &todoTasks)
}

// UpdateTodoTask is a handler function for updating a todo task.
func (h *MyDatabase) UpdateTodoTask(c *gin.Context) {
	// Extract the ID parameter from the request URL.
	id := c.Param("id")

	// Declare a variable to store the request body.
	body := TodoTaskRequestBody{}

	// Receive request body.
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Declare a variable to store the retrieved todo task.
	var todoTask models.TodoTask

	// Retrieve the todo task from the database by its ID.
	if result := h.DB.First(&todoTask, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	// Update the todo task's fields with the values from the request body.
	todoTask.Title = body.Title
	todoTask.Description = body.Description
	todoTask.State = body.State

	// Save the updated todo task to the database.
	h.DB.Save(&todoTask)

	// Respond with the updated todo task.
	c.JSON(http.StatusOK, &todoTask)
}
