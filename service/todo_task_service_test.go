package service

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
