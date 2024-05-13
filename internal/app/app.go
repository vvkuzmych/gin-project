package app

import (
	"github.com/vkuzmich/gin-project/repository"
	"github.com/vkuzmich/gin-project/service"
	"gorm.io/gorm"
)

var _ Interface = (*App)(nil)

type Interface interface {
	TodoTaskRepository() repository.TodoTaskRepository
	TodoTaskService() service.TodoTaskService
}

type App struct {
	todoTaskService service.TodoTaskService

	todoTaskRepository repository.TodoTaskRepository
}

//func (a *App) TodoTaskRepository() repository.TodoTaskRepository {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (a *App) TodoTaskService() service.TodoTaskService {
//	//TODO implement me
//	panic("implement me")
//}

func (a *App) TodoTaskRepository() repository.TodoTaskRepository {
	return a.todoTaskRepository
}

func (a *App) TodoTaskService() service.TodoTaskService {
	return a.todoTaskService
}

func Build(db *gorm.DB) *App {

	var (
		todoTaskRepository = repository.NewTodoTaskRepository(db)
		todoTaskService    = service.NewTodoTaskService(todoTaskRepository)
	)

	app := &App{
		todoTaskRepository: todoTaskRepository,
		todoTaskService:    todoTaskService,
	}

	return app
}
