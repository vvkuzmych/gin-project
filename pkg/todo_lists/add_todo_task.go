package todo_lists

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/vkuzmich/gin-project/pkg/common/models"
)

// AddTodoTaskRequestBody represents the structure of the request body
// expected when adding a new todo task.
type AddTodoTaskRequestBody struct {
    Title       string `json:"title"`       // Title of the todo task
    Description string `json:"description"` // Description of the todo task
    State       bool   `json:"state"`       // State of the todo task (completed or not)
}

// AddTodoTask is a handler function for adding a new todo task.
func (h handler) AddTodoTask(c *gin.Context) {
    body := AddTodoTaskRequestBody{}

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
