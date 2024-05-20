package repository

import (
	"gorm.io/gorm"
	"log"
	"testing"
)

var testDB *gorm.DB
var todoTaskRepo TodoTaskRepository

func TestMain(m *testing.M) {

	log.Println("Setting up the postgres DB for mocking.......")
	ctx, container, db := InitiateContainerCreation()

	testDB = db

	todoTaskRepo = NewTodoTaskRepository(testDB)

	defer func() {
		if container != nil && container.IsRunning() {
			log.Println("Shutting down the phanes postgres db test container....")
			container.Terminate(ctx)
		}
	}()
	m.Run()
}
