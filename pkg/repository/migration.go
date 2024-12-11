package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"

	"github.com/golang-migrate/migrate/v4"
	scripts "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Logger struct {
	enabled bool
	multi   bool
}

// repository represents a database connection.
type repository struct {
	db *gorm.DB // db is a reference to the database connection
}

func (t Logger) Printf(format string, v ...interface{}) {
	if t.enabled {
		log.Printf("migrate: %s", fmt.Sprintf(format, v...))
	}
}

func (t Logger) Verbose() bool {
	return t.multi
}

func SetupConnection() (*gorm.DB, error) {
	var db *gorm.DB
	log.Println("Setting up database connection...")
	domain := getDomain()
	db, err := GetConnection(domain)

	if err != nil {
		return nil, err
	}

	URL := "file://./db/migration/"
	err = ExecMigration(domain, URL, true, true)
	if err != nil {
		return nil, err
	}
	return db, err
}

func GetConnection(domain string) (*gorm.DB, error) {
	config := &gorm.Config{}

	db, err := gorm.Open(postgres.Open(domain), config)
	if err != nil {
		return nil, fmt.Errorf("unable to open db: %s", err)
	}

	return db, nil
}

func getDomain() string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Get values from Viper
	dbUserName := viper.GetString("POSTGRES_USER")
	dbPassword := viper.GetString("POSTGRES_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbName := viper.GetString("DB_NAME")

	domain := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUserName, dbName, dbPassword)
	return domain
}

// method is used to execute the migration up scripts to the
// phanes database using golang-migrate
func ExecMigration(connection string, url string, multi bool, loggerEnabled bool) error {
	db, err := sql.Open("postgres", connection)
	defer func() {
		err = db.Close()
		if err != nil {
			log.Println("unable to close DB", err)
		}
	}()

	if err != nil {
		log.Println("unable to open DB", err)
		return err
	}

	driver, err := scripts.WithInstance(db, &scripts.Config{})
	if err != nil {
		log.Println("Error while loading driver", err)
		return err
	}

	// Load and execute the migration scripts using golang-migrate lib
	m, err := migrate.NewWithDatabaseInstance(
		url,
		"postgres", driver)

	if err != nil {
		log.Println("Error initializing migrations", err)
		return err
	}

	// log all migrations run:
	m.Log = Logger{multi: multi, enabled: loggerEnabled}

	if err != nil {
		log.Println("Error while loading migration scripts", err)
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println("Error while running migration up scripts", err)
		return err
	}
	return nil
}
