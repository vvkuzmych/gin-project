package repository

import (
	"errors"
	"github.com/vkuzmich/gin-project/contextLogger"
	"github.com/vkuzmich/gin-project/internal/model"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TodoTaskRepository interface {
	CreateTodoTask(ctx context.Context, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error)
	DeleteTodoTask(ctx context.Context, id string) error
	GetTodoTask(ctx context.Context, id string) (model.TodoTask, error)
	GetAllTodoTasks(ctx context.Context) ([]model.TodoTask, error)
	UpdateTodoTask(ctx context.Context, id string, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error)
}

func NewTodoTaskRepository(db *gorm.DB) TodoTaskRepository {
	return repository{db}
}

// CreateTodoTask is a handler function for adding a new.
func (r repository) CreateTodoTask(ctx context.Context, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)

	// Validate the todoTaskPayload
	if err := todoTaskPayload.ValidateTodoTaskPayload(); err != nil {
		return model.TodoTask{}, err
	}

	todoTask := model.TodoTask{
		Title:       todoTaskPayload.Title,
		Description: todoTaskPayload.Description,
		State:       todoTaskPayload.State,
	}

	if r.db == nil {
		return model.TodoTask{}, errors.New("database connection is nil")
	}

	result := r.db.Create(&todoTask)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("error while creating todo_task")
		return model.TodoTask{}, result.Error
	}
	logger.Info().Msg("TodoTask created")
	// Return the created TodoTask with the generated ID
	return todoTask, nil
}

func (r repository) DeleteTodoTask(ctx context.Context, id string) error {
	logger := contextLogger.ContextLog(ctx)

	if err := r.db.Delete(&model.TodoTask{}, id).Error; err != nil {
		logger.Error().Err(err).Msg("error while deleting todo_task")
		return errors.New("Invalid id")
	}

	logger.Info().Msg("TodoTask deleted")
	return nil
}

func (r repository) GetTodoTask(ctx context.Context, id string) (model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)

	if id == "" {
		logger.Info().Str("todo_task_id", id).Msg("invalid todo_task_id")
		return model.TodoTask{}, errors.New("Invalid id")
	}

	var todoTask model.TodoTask
	result := r.db.Where("id = ?", id).First(&todoTask)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("error while getting todo_task")
		return model.TodoTask{}, result.Error
	}

	logger.Info().Msg("TodoTask was found")
	// Return single TodoTask
	return todoTask, nil
}

func (r repository) GetAllTodoTasks(ctx context.Context) ([]model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)

	var todoTask []model.TodoTask
	if err := r.db.Find(&todoTask).Error; err != nil {
		logger.Error().Err(err).Msg("error while fetching todo_tasks")
		return []model.TodoTask{}, err
	}

	logger.Info().Msg("Get all TodoTasks")
	// Return array of TodoTasks
	return todoTask, nil
}

func (r repository) UpdateTodoTask(ctx context.Context, id string, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)

	var todoTask model.TodoTask

	if err := r.db.Find(&todoTask, id).Error; err != nil {
		logger.Error().Err(err).Msg("error while getting todo_task")
		return model.TodoTask{}, err
	}

	todoTask.Title = todoTaskPayload.Title
	todoTask.Description = todoTaskPayload.Description
	todoTask.State = todoTaskPayload.State

	result := r.db.Model(model.TodoTask{}).Where("id = ?", id).Updates(&todoTask)
	if result.Error != nil {
		logger.Error().Err(result.Error).Str("todo_task_id", id).Msgf("Error while updating todo_task")
		return model.TodoTask{}, result.Error
	}
	logger.Info().Msg("Get updated TodoTask")
	// Return the created TodoTask with the generated ID
	return todoTask, nil
}
