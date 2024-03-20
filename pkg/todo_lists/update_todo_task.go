package todo_lists

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/vkuzmich/gin-project/pkg/common/models"
)

// UpdateTodoTaskRequestBody represents the structure of the request body
// expected when updating a todo task.
type UpdateTodoTaskRequestBody struct {
    Title       string `json:"title"`       // Title of the todo task
    Description string `json:"description"` // Description of the todo task
    State       bool   `json:"state"`       // State of the todo task (completed or not)
}

// UpdateTodoTask is a handler function for updating a todo task.
func (h handler) UpdateTodoTask(c *gin.Context) {
    // Extract the ID parameter from the request URL.
    id := c.Param("id")

    // Declare a variable to store the request body.
    body := UpdateTodoTaskRequestBody{}

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
