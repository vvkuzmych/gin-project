package service

import (
	"github.com/gin-gonic/gin"
	"github.com/vkuzmich/gin-project/contextLogger"
	"github.com/vkuzmich/gin-project/internal/model"
	"github.com/vkuzmich/gin-project/repository"
)

// repository represents a database connection.
type TodoTaskService interface {
	AddTodoTask(c *gin.Context, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error)
	DeleteTodoTask(c *gin.Context, id string) error
	GetTodoTask(ctx *gin.Context, id string) (model.TodoTask, error)
	GetTodoTasks(ctx *gin.Context) ([]model.TodoTask, error)
	UpdateTodoTask(ctx *gin.Context, id string, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error)
}

func NewTodoTaskService(todoTaskRepository repository.TodoTaskRepository) TodoTaskService {
	return todoTaskService{
		todoTaskRepository,
	}
}

type todoTaskService struct {
	todoTaskRepository repository.TodoTaskRepository
}

// AddTodoTask is a handler function for adding a new todoTask.
func (s todoTaskService) AddTodoTask(ctx *gin.Context, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)

	todoTask, err := s.todoTaskRepository.CreateTodoTask(ctx, todoTaskPayload)

	if err != nil {
		logger.Error().Err(err).Msg("Fail to create todo_task")
		return model.TodoTask{}, err
	}
	logger.Info().Msg("Successfully created todo_task")
	return todoTask, nil
}

func (s todoTaskService) DeleteTodoTask(ctx *gin.Context, id string) error {
	logger := contextLogger.ContextLog(ctx)

	err := s.todoTaskRepository.DeleteTodoTask(ctx, id)

	if err != nil {
		logger.Error().Err(err).Msg("Fail to delete todo_task")
		return err
	}
	logger.Info().Msg("Successfully delete todo_task")
	return nil
}

func (s todoTaskService) GetTodoTask(ctx *gin.Context, id string) (model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)
	todoTask, err := s.todoTaskRepository.GetTodoTask(ctx, id)

	if err != nil {
		logger.Error().Err(err).Msg("Fail to get todo_task")
		return model.TodoTask{}, err
	}
	logger.Info().Msg("Successfully get todo_task")
	return todoTask, nil
}

func (s todoTaskService) GetTodoTasks(ctx *gin.Context) ([]model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)
	todoTask, err := s.todoTaskRepository.GetAllTodoTasks(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Fail to get todo_tasks")
		return []model.TodoTask{}, err
	}
	logger.Info().Msg("Successfully get todo_tasks")
	return todoTask, nil
}

func (s todoTaskService) UpdateTodoTask(ctx *gin.Context, id string, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)
	todoTask, err := s.todoTaskRepository.UpdateTodoTask(ctx, id, todoTaskPayload)

	if err != nil {
		logger.Error().Err(err).Msg("Fail to get todo_task")
		return model.TodoTask{}, err
	}
	logger.Info().Msg("Successfully get todo_task")

	return todoTask, nil
}
