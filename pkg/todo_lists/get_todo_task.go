package todo_lists

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/vkuzmich/gin-project/pkg/common/models"
)

// GetTodoTask is a handler function for retrieving a single todo task by ID.
func (h handler) GetTodoTask(c *gin.Context) {
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
