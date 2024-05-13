package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vkuzmich/gin-project/contextLogger"
	"github.com/vkuzmich/gin-project/internal/model"
	"gorm.io/gorm"
	"net/http"

	//"github.com/go-playground/validator/v10"
	"github.com/vkuzmich/gin-project/service"
)

func RegisterTodoTaskHandlers(
	r *gin.RouterGroup,
	todoTaskService service.TodoTaskService,
) {

	res := TodoTaskResource{todoTaskService}

	todoTask := r.Group("/todo_tasks")
	{
		todoTask.POST("/", res.AddTodoTaskRoute)
		todoTask.GET("/", res.GetTodoTasksRoute)
		todoTask.GET("/:id", res.GetTodoTaskRoute)
		todoTask.PUT("/:id", res.UpdateTodoTaskRoute)
		todoTask.DELETE("/:id", res.DeleteTodoTaskRoute)
	}
}

type TodoTaskResource struct {
	todoTaskService service.TodoTaskService
}

// AddTodoTaskRequestBody represents the structure of the request body
// expected when adding a new todo task.
type TodoTaskRequestBody struct {
	Title       string `json:"title"`       // Title of the todo task
	Description string `json:"description"` // Description of the todo task
	State       bool   `json:"state"`       // State of the todo task (completed or not)
}

func (r TodoTaskResource) AddTodoTaskRoute(ctx *gin.Context) {
	logger := contextLogger.ContextLog(ctx)
	logger.Info().Msg("AddTodoTask endpoint hit")

	body := TodoTaskRequestBody{}
	// Receive request body
	if err := ctx.BindJSON(&body); err != nil {
		logger.Error().Err(err).Msg("Error in Binding todo_task payload from request")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Assign values from request body to todoTask
	todoTask := model.TodoTaskPayload{
		Title:       body.Title,
		Description: body.Description,
		State:       body.State,
	}

	// Create todoTask in the database
	result, err := r.todoTaskService.AddTodoTask(ctx, &todoTask)
	if err != nil {
		logger.Error().Err(err).Msg("Error in processing todo_task")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logger.Info().Msg("AddTodoTask endpoint successfully created todo_task")
	// Respond with the created todo_task
	ctx.JSON(http.StatusOK, &result)
}

func (r TodoTaskResource) GetTodoTasksRoute(ctx *gin.Context) {
	logger := contextLogger.ContextLog(ctx)
	logger.Info().Msg("GetTodoTasks endpoint hit")
	// Retrieve todo_tasks from the database.
	todoTasks, err := r.todoTaskService.GetTodoTasks(ctx)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			logger.Info().Msg("no todo_tasks")
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		logger.Error().Err(err).Msg("Error in getting todo_tasks")
		// Abort the request with an error if retrieval fails.
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Respond with the retrieved todo tasks.
	ctx.JSON(http.StatusOK, &todoTasks)
}

func (r TodoTaskResource) GetTodoTaskRoute(ctx *gin.Context) {
	logger := contextLogger.ContextLog(ctx)
	logger.Info().Msg("GetTodoTask endpoint hit")
	// Extract the ID parameter from the request URL.
	id := ctx.Param("id")

	// Retrieve the todo_task from the database by its ID.
	todoTask, err := r.todoTaskService.GetTodoTask(ctx, id)
	if err != nil {
		// Abort the request with an error if retrieval fails.
		logger.Error().Err(err).Str("todo_task_id", id).Msg("Error in getting todo_task")
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Respond with a success status.
	logger.Info().Msg("GetTodoTask endpoint successfully get todo_task")
	ctx.JSON(http.StatusOK, &todoTask)
}

func (r TodoTaskResource) UpdateTodoTaskRoute(ctx *gin.Context) {
	// Extract the ID parameter from the request URL.
	id := ctx.Param("id")

	// Declare a variable to store the request body.
	body := TodoTaskRequestBody{}

	// Receive request body.
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Assign values from request body to todoTask
	todoTaskPayload := model.TodoTaskPayload{
		Title:       body.Title,
		Description: body.Description,
		State:       body.State,
	}

	// Retrieve the todo_task from the database by its ID.
	todoTask, err := r.todoTaskService.UpdateTodoTask(ctx, id, &todoTaskPayload)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Respond with the updated todo_task.
	ctx.JSON(http.StatusOK, &todoTask)
}

func (r TodoTaskResource) DeleteTodoTaskRoute(ctx *gin.Context) {
	logger := contextLogger.ContextLog(ctx)
	logger.Info().Msg("DeleteTodoTask endpoint hit")
	// Extract the ID parameter from the request URL.
	id := ctx.Param("id")

	// Retrieve the todo_task from the database by its ID.
	err := r.todoTaskService.DeleteTodoTask(ctx, id)
	if err != nil {
		// Abort the request with an error if retrieval fails.
		logger.Error().Err(err).Str("todo_task_id", id).Msg("Error in deleting todo_task")
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Respond with a success status.
	ctx.Status(http.StatusOK)
}
