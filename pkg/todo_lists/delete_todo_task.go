package todo_lists

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/vkuzmich/gin-project/pkg/common/models"
)

// DeleteTodoTask is a handler function for deleting a todo task.
func (h handler) DeleteTodoTask(c *gin.Context) {
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
