package todo_lists

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/vkuzmich/gin-project/pkg/common/models"
)

// GetTodoTasks is a handler function for retrieving todo tasks.
func (h handler) GetTodoTasks(c *gin.Context) {
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
