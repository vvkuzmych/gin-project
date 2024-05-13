package repository

import (
	"github.com/vkuzmich/gin-project/contextLogger"
	"github.com/vkuzmich/gin-project/internal/model"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

// repository represents a database connection.
type repository struct {
	db *gorm.DB // db is a reference to the database connection
}

type TodoTaskRepository interface {
	CreateTodoTask(ctx context.Context, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error)
	DeleteTodoTask(ctx context.Context, id string) error
	GetTodoTask(ctx context.Context, id string) (model.TodoTask, error)
	GetTodoTasks(ctx context.Context) ([]model.TodoTask, error)
	UpdateTodoTask(ctx context.Context, id string, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error)
}

func NewTodoTaskRepository(db *gorm.DB) TodoTaskRepository {
	return repository{db}
}

// CreateTodoTask is a handler function for adding a new.
func (r repository) CreateTodoTask(ctx context.Context, todoTaskPayload *model.TodoTaskPayload) (model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)

	todoTask := model.TodoTask{
		Title:       todoTaskPayload.Title,
		Description: todoTaskPayload.Description,
		State:       todoTaskPayload.State,
	}

	if err := r.db.Create(&todoTask).Error; err != nil {
		logger.Error().Err(err).Msg("error while creating todo_task")
		return model.TodoTask{}, err
	}
	logger.Info().Msg("TodoTask created")
	// Return the created TodoTask with the generated ID
	return todoTask, nil
}

func (r repository) DeleteTodoTask(ctx context.Context, id string) error {
	logger := contextLogger.ContextLog(ctx)

	if err := r.db.Delete(&model.TodoTask{}, id).Error; err != nil {
		logger.Error().Err(err).Msg("error while deleting todo_task")
		return err
	}

	logger.Info().Msg("TodoTask deleted")
	// Return the created TodoTask with the generated ID
	return nil
}

func (r repository) GetTodoTask(ctx context.Context, id string) (model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)

	var todoTask model.TodoTask
	if err := r.db.First(&todoTask, id).Error; err != nil {
		logger.Error().Err(err).Msg("error while getting todo_task")
		return model.TodoTask{}, err
	}

	logger.Info().Msg("TodoTask was found")
	// Return the created TodoTask with the generated ID
	return todoTask, nil
}

func (r repository) GetTodoTasks(ctx context.Context) ([]model.TodoTask, error) {
	logger := contextLogger.ContextLog(ctx)

	var todoTask []model.TodoTask
	if err := r.db.Find(&todoTask).Error; err != nil {
		logger.Error().Err(err).Msg("error while getting todo_task")
		return []model.TodoTask{}, err
	}

	logger.Info().Msg("Get all TodoTasks")
	// Return the created TodoTask with the generated ID
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
