package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vkuzmich/gin-project/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func GetMockedDBInstance() (*gorm.DB, sqlmock.Sqlmock) {
	var mockedDB *gorm.DB
	conn, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil || conn == nil {
		panic(fmt.Sprintf("Failed to open mock sql db, got error: %v", err))
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 conn,
		PreferSimpleProtocol: true,
	})

	if db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true}); err != nil || db == nil {
		panic(fmt.Sprintf("Failed to open gorm v2 db, got error: %v", err))
	} else {
		mockedDB = db
	}
	return mockedDB, mock
}

func TestCreateTodoTask(t *testing.T) {
	type world struct {
		todoTaskPayload model.TodoTaskPayload
	}
	tests := []struct {
		name            string
		ctx             context.Context
		todoTaskPayload *model.TodoTaskPayload
		expectedError   error
		setup           func(t *testing.T, d *world)
	}{
		{
			name:          "Valid Payload",
			ctx:           context.Background(),
			expectedError: nil,
			setup: func(t *testing.T, d *world) {
				d.todoTaskPayload.State = true
			},
		},
		{
			name:          "invalid Payload title",
			ctx:           context.Background(),
			expectedError: errors.New("validation fails: Key: 'TodoTaskPayload.Title' Error:Field validation for 'Title' failed on the 'required' tag"),
			setup: func(t *testing.T, d *world) {
				d.todoTaskPayload.Title = ""
				d.todoTaskPayload.State = true
			},
		},
		{
			name:          "invalid Payload description",
			ctx:           context.Background(),
			expectedError: errors.New("validation fails: Key: 'TodoTaskPayload.Description' Error:Field validation for 'Description' failed on the 'required' tag"),
			setup: func(t *testing.T, d *world) {
				d.todoTaskPayload.Description = ""
				d.todoTaskPayload.State = true
			},
		},
		{
			name:          "invalid Payload state",
			ctx:           context.Background(),
			expectedError: errors.New("validation fails: Key: 'TodoTaskPayload.State' Error:Field validation for 'State' failed on the 'required' tag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &world{
				todoTaskPayload: model.TodoTaskPayload{
					Title:       "Test Task",
					Description: "Test Description",
				},
			}

			if tt.setup != nil {
				tt.setup(t, d)
			}
			// Create a repository instance with the mocked database
			repo := repository{db: testDB}

			// Call the function with the test payload and context
			_, err := repo.CreateTodoTask(tt.ctx, &d.todoTaskPayload)

			//Check the error returned
			assert.Equal(t, tt.expectedError, err)

			t.Cleanup(func() {
				AfterEach()
			})
		})
	}
}

func TestCreateTodoTaskError(t *testing.T) {
	mockedDB, mock := GetMockedDBInstance()
	mainDB := testDB
	testDB = mockedDB
	todoTaskPayload := model.TodoTaskPayload{
		Title:       "Test Task",
		Description: "Test Description",
		State:       true,
	}
	var todoTaskRepository = NewTodoTaskRepository(testDB)
	mock.ExpectExec(`INSERT INTO "todo_tasks" `).WithArgs(&todoTaskPayload).WillReturnError(errors.New("error"))

	result, err := todoTaskRepository.CreateTodoTask(context.Background(), &todoTaskPayload)
	testDB = mainDB
	assert.NotNil(t, result)
	assert.Error(t, err, "error")
}

func TestDeleteTodoTask(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		id            string
		expectedError error
	}{
		{
			name:          "Valid ID",
			ctx:           context.Background(),
			id:            "1",
			expectedError: nil,
		},
		{
			name:          "Invalid ID",
			ctx:           context.Background(),
			id:            "", // Invalid ID
			expectedError: errors.New("Invalid id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository{db: testDB}

			todoTaskPayload := model.TodoTaskPayload{
				Title:       "Test Task",
				Description: "Test Description",
				State:       true,
			}

			// Create a repository instance with the mocked database
			_, err := CreateTodoTask(t, repo, todoTaskPayload)
			assert.NoError(t, err)

			// Call the function with the test context and ID
			err = repo.DeleteTodoTask(tt.ctx, tt.id)

			// Check for any errors
			assert.Equal(t, tt.expectedError, err)
		})
		t.Cleanup(func() {
			AfterEach()
		})
	}
}

func TestGetTodoTask(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		id             string
		expectedError  error
		expectedResult model.TodoTask
	}{
		{
			name:          "Valid ID",
			ctx:           context.Background(),
			id:            "1",
			expectedError: nil,
			expectedResult: model.TodoTask{
				Title:       "Test Task",
				Description: "Test Description",
				State:       true,
			},
		},
		{
			name:           "Invalid ID",
			ctx:            context.Background(),
			id:             "", // Invalid ID
			expectedError:  errors.New("Invalid id"),
			expectedResult: model.TodoTask{},
		},
		{
			name:           "Invalid ID",
			ctx:            context.Background(),
			id:             "40", // Invalid ID
			expectedError:  errors.New("record not found"),
			expectedResult: model.TodoTask{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository{db: testDB}

			todoTaskPayload := model.TodoTaskPayload{
				Title:       "Test Task",
				Description: "Test Description",
				State:       true,
			}

			// Create a repository instance with the mocked database
			_, err := CreateTodoTask(t, repo, todoTaskPayload)
			assert.NoError(t, err)

			// Call the function with the test context and ID
			result, resultErr := repo.GetTodoTask(tt.ctx, tt.id)

			// Check for any errors
			assert.Equal(t, tt.expectedError, resultErr)
			assert.Equal(t, tt.expectedResult.Title, result.Title)
			assert.Equal(t, tt.expectedResult.Description, result.Description)
			assert.Equal(t, tt.expectedResult.State, result.State)
		})
		t.Cleanup(func() {
			AfterEach()
		})
	}
}

func CreateTodoTask(t *testing.T, repo repository, todoTaskPayload model.TodoTaskPayload) (uint, error) {
	createdTodoTask, err := repo.CreateTodoTask(context.Background(), &todoTaskPayload)
	assert.NoError(t, err)
	return createdTodoTask.ID, nil
}
