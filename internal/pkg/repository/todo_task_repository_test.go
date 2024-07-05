//go:build UnitTest
// +build UnitTest

package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vkuzmich/gin-project/internal/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
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

func TestGetListTodoTasks(t *testing.T) {
	repo := repository{db: testDB}
	CreateTodoTasksList(t, repo)

	tests := []struct {
		name           string
		ctx            context.Context
		expectedError  error
		expectedResult []model.TodoTask
	}{
		{
			name:          "Successfully get list of todo_tasks",
			ctx:           context.Background(),
			expectedError: nil,
			expectedResult: []model.TodoTask{
				{
					Title:       "Test Task 1",
					Description: "Test Description 1",
					State:       true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call the function with the test context
			result, resultErr := repo.GetTodoTasks(tt.ctx)

			// Check for any errors
			assert.Equal(t, tt.expectedError, resultErr)
			assert.Equal(t, tt.expectedResult[0].Title, result[0].Title)
			assert.Equal(t, tt.expectedResult[0].Description, result[0].Description)
			assert.Equal(t, tt.expectedResult[0].State, result[0].State)
			assert.Equal(t, 3, len(result))
		})
		t.Cleanup(func() {
			AfterEach()
		})
	}
}

func TestGetAllTodoTasksError(t *testing.T) {
	mockedDB, mock := GetMockedDBInstance()
	mainDB := testDB
	testDB = mockedDB
	var mockTodoTaskRepository = NewTodoTaskRepository(testDB)

	mock.ExpectQuery(`SELECT * FROM "todo_tasks"`).WillReturnError(fmt.Errorf("error while fetching todo_tasks"))

	_, err := mockTodoTaskRepository.GetTodoTasks(context.Background())
	testDB = mainDB
	assert.NotNil(t, err)
}

func TestUpdateTodoTask(t *testing.T) {
	type world struct {
		todoTaskPayload model.TodoTaskPayload
	}
	tests := []struct {
		name            string
		ctx             context.Context
		todoTaskPayload *model.TodoTaskPayload
		id              string
		expectedError   error
		setup           func(t *testing.T, d *world)
	}{
		{
			name:          "Successful update",
			ctx:           context.Background(),
			expectedError: nil,
			id:            "1",
			setup: func(t *testing.T, d *world) {
				d.todoTaskPayload.State = true
			},
		},
		{
			name:          "invalid id",
			ctx:           context.Background(),
			expectedError: errors.New("Invalid id"),
			id:            "",
			setup: func(t *testing.T, d *world) {
				d.todoTaskPayload.State = true
			},
		},
		{
			name:          "invalid Payload title",
			ctx:           context.Background(),
			id:            "1",
			expectedError: errors.New("validation fails: Key: 'TodoTask.Title' Error:Field validation for 'Title' failed on the 'required' tag"),
			setup: func(t *testing.T, d *world) {
				d.todoTaskPayload.Title = ""
				d.todoTaskPayload.State = true
			},
		},
		{
			name:          "invalid Payload description",
			ctx:           context.Background(),
			id:            "1",
			expectedError: errors.New("validation fails: Key: 'TodoTask.Description' Error:Field validation for 'Description' failed on the 'required' tag"),
			setup: func(t *testing.T, d *world) {
				d.todoTaskPayload.Description = ""
				d.todoTaskPayload.State = true
			},
		},
		{
			name:          "invalid Payload state",
			ctx:           context.Background(),
			id:            "1",
			expectedError: errors.New("validation fails: Key: 'TodoTask.State' Error:Field validation for 'State' failed on the 'required' tag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &world{
				todoTaskPayload: model.TodoTaskPayload{
					Title:       "New Task",
					Description: "New Description",
				},
			}

			if tt.setup != nil {
				tt.setup(t, d)
			}
			// Create a repository instance with the mocked database
			repo := repository{db: testDB}
			CreateTodoTasksList(t, repo)
			// Call the function with the test payload and context
			result, err := repo.UpdateTodoTask(tt.ctx, tt.id, &d.todoTaskPayload)

			//Check the error returned
			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.Equal(t, d.todoTaskPayload.Title, result.Title)
				assert.Equal(t, d.todoTaskPayload.Description, result.Description)
				assert.Equal(t, d.todoTaskPayload.State, result.State)
			}

			t.Cleanup(func() {
				AfterEach()
			})
		})
	}
}

func CreateTodoTasksList(t *testing.T, repo repository) {
	// Create a repository instance with the mocked database
	for i := 1; i <= 3; i++ {
		todoTaskPayload := model.TodoTaskPayload{
			Title:       "Test Task " + strconv.Itoa(i),
			Description: "Test Description " + strconv.Itoa(i),
			State:       true,
		}
		_, err := CreateTodoTask(t, repo, todoTaskPayload)
		assert.NoError(t, err)
	}
}

func CreateTodoTask(t *testing.T, repo repository, todoTaskPayload model.TodoTaskPayload) (uint, error) {
	createdTodoTask, err := repo.CreateTodoTask(context.Background(), &todoTaskPayload)
	assert.NoError(t, err)
	return createdTodoTask.ID, nil
}
