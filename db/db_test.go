package db

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockDB is a mock implementation of *gorm.DB
type MockDB struct {
	mock.Mock
}

// AutoMigrate is a mock method that simulates the behavior of the AutoMigrate method of *gorm.DB
func (m *MockDB) AutoMigrate(dst ...interface{}) error {
	args := m.Called(dst)
	return args.Error(0)
}

func TestInit(t *testing.T) {
	tests := []struct {
		name           string
		mockDB         *MockDB
		autoMigrations func(db *gorm.DB) error
		connectingFunc func(url string) (*gorm.DB, error)
		expectedError  string
	}{
		{
			name:   "Success",
			mockDB: new(MockDB),
			autoMigrations: func(db *gorm.DB) error {
				return nil
			},
			connectingFunc: func(url string) (*gorm.DB, error) {
				return mockToGormDB(new(MockDB)), nil
			},
			expectedError: "",
		},
		{
			name:   "FailureAutoMigrations",
			mockDB: new(MockDB),
			autoMigrations: func(db *gorm.DB) error {
				return errors.New("fails to migrate")
			},
			connectingFunc: func(url string) (*gorm.DB, error) {
				return mockToGormDB(new(MockDB)), nil
			},
			expectedError: "fails to migrate",
		},
		{
			name:   "FailureConnectingToDB",
			mockDB: new(MockDB),
			autoMigrations: func(db *gorm.DB) error {
				return nil
			},
			connectingFunc: func(url string) (*gorm.DB, error) {
				return nil, errors.New("fails to connect")
			},
			expectedError: "fails to connect",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ConnectingToDB = tc.connectingFunc
			AutoMigrations = tc.autoMigrations

			db, err := Init("mock-db-url")

			if tc.expectedError != "" {
				assert.Error(t, err, tc.expectedError)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}

// Function to convert MockDB to *gorm.DB
func mockToGormDB(mockDB *MockDB) *gorm.DB {
	// Implement a simple adapter to satisfy *gorm.DB interface
	return &gorm.DB{}
}

func TestConnectionToDB(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		wantError     bool
		expectedError string
	}{
		{
			name:          "Success",
			url:           "postgres://admin:admin123@localhost:5432/gin_pron",
			wantError:     false,
			expectedError: "",
		},
		//{
		//	name:          "Failure",
		//	url:           "invalid-db-url",
		//	wantError:     true,
		//	expectedError: "error connecting to database",
		//},
	}

	// Mock the ConnectionToDB function to return an error
	ConnectingToDB = func(url string) (*gorm.DB, error) {
		// Return a custom error message for the specified URL
		if url == "invalid-db-url" {
			return nil, fmt.Errorf("failed to parse as DSN (invalid dsn)")
		}
		// Return a nil DB and nil error for other URLs
		return nil, nil
	}
	defer func() { ConnectingToDB = ConnectionToDB }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := ConnectionToDB(tt.url)
			if tt.wantError {
				assert.Error(t, err, tt.expectedError)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)

				// Close the database connection
				sqlDB, err := db.DB()
				assert.NoError(t, err)
				assert.NotNil(t, sqlDB)
				defer sqlDB.Close()
			}
		})
	}
}

func TestAutoMigration(t *testing.T) {
	tests := []struct {
		name      string
		db        *gorm.DB
		wantError bool
		errMsg    error
	}{
		{
			name:      "Success",
			db:        createSQLiteDB(t),
			wantError: false,
			errMsg:    nil,
		},
		{
			name:      "NilDB",
			db:        nil,
			wantError: true,
			errMsg:    errors.New("nil database connection"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AutoMigration(tt.db)
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create an SQLite in-memory database for testing
func createSQLiteDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	return db
}
