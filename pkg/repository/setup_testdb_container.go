package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/vkuzmich/gin-project/pkg/db"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	migrationScripts "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "gorm.io/driver/postgres"
)

const (
	DB_USERNAME         = "testuser"
	DB_PASSWORD         = "testpwd"
	DB_NAME             = "db-test"
	DB_PORT             = "5432"
	DB_PROVIDER         = "postgres"
	DB_MIGRATION_SOURCE = "file://./../db/migration/"
)

type postgresContainer struct {
	testcontainers.Container
	host string
	port string
}

var connStatement string

type InitiateContainerOption int

const (
	SkipRunningMigrations InitiateContainerOption = iota + 1
)

func InitiateContainerCreation(options ...InitiateContainerOption) (context.Context, *postgresContainer, *gorm.DB) {
	ctx, container, connSt, err := CreatePostgresContainer()
	connStatement = connSt
	if err != nil {
		panic(fmt.Sprintf("unable to setup the postgres test container: %s", err))
	}

	var (
		runMigration = true
	)

	for _, o := range options {
		if o == SkipRunningMigrations {
			runMigration = false
		}
	}

	// Once the test DB container is created we need to execute the
	// migration scripts
	log.Println("Connection statement->", connStatement)
	if runMigration {
		err = ExecMigration(connStatement, DB_MIGRATION_SOURCE, true, true)
		if err != nil {
			panic(fmt.Sprintf("failure running migrations: %s", err))
		}
	}
	// ..
	var db *gorm.DB
	db, err = GetConnection(connStatement)
	if err != nil {
		panic(fmt.Sprintf("unable to connect to the postgres test container: %s", err))
	}
	return ctx, container, db
}

// method is used to create and configure postgres container using test container
// for mocking unit tests for the phanes. Test container manages life cycle for phanes
// test DB container. The return type of the method is the SQL connection string.
func CreatePostgresContainer() (context.Context, *postgresContainer, string, error) {
	ctx := context.Background()
	dbname := DB_NAME
	user := DB_USERNAME
	password := DB_PASSWORD
	port, err := nat.NewPort("tcp", DB_PORT)
	if err != nil {
		log.Fatal("Error creating new Network Port", err)
	}
	containerMappedPort := port.Port()
	container, err := setupPostgres(ctx,
		WithPort(containerMappedPort),
		WithInitialDatabase(user, password, dbname),
		WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)

	if err != nil {
		log.Printf("unable to create postgres image: %v\n", err)
		return nil, nil, "", err
	}

	containerPort, err := container.MappedPort(ctx, port)
	if err != nil {
		log.Println("unable to map the port for database in the test container", err)
		return nil, nil, "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Println("unable to assign the host for database in the test container", err)
		return nil, nil, "", err
	}
	connStr := DB_PROVIDER + "://" + user + ":" + password + "@" + host + ":" + containerPort.Port() + "/" + dbname + "?sslmode=disable"

	container.host = host
	container.port = containerPort.Port()

	return ctx, container, connStr, err
}

func (c *postgresContainer) GetHost() string {
	return c.host
}

func (c *postgresContainer) GetPort() string {
	return c.port
}

type postgresContainerOption func(req *testcontainers.ContainerRequest)

func WithWaitStrategy(strategies ...wait.Strategy) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.WaitingFor = wait.ForAll(strategies...).WithDeadline(1 * time.Minute)
	}
}

// method is used to configure Database port for the test container
func WithPort(port string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.ExposedPorts = append(req.ExposedPorts, port)
	}
}

// method is used to configure Database userName, password and Database name for the test container
func WithInitialDatabase(user string, password string, dbName string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.Env["POSTGRES_USER"] = user
		req.Env["POSTGRES_PASSWORD"] = password
		req.Env["POSTGRES_DB"] = dbName
	}
}

// setupPostgres creates an instance of the postgres container type
func setupPostgres(ctx context.Context, opts ...postgresContainerOption) (*postgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13.8",
		Env:          map[string]string{},
		ExposedPorts: []string{},
		Cmd:          []string{"postgres", "-c", "fsync=off", "-N", "500"},
		AutoRemove:   true,
	}

	for _, opt := range opts {
		opt(&req)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &postgresContainer{Container: container}, nil
}

func AfterEach() {
	db, err := sql.Open("postgres", connStatement)
	if err != nil {
		log.Fatal("Error connecting to postgres AfterEach", err)
	}

	defer db.Close()

	if err != nil {
		log.Println("unable to open DB connection using connection statement..", err)
	}

	driver, err := migrationScripts.WithInstance(db, &migrationScripts.Config{})
	if err != nil {
		log.Println("Error while loading driver config", err)
	}

	// Load and execute the migration scripts using golang-migrate lib
	m, err := migrate.NewWithDatabaseInstance(
		DB_MIGRATION_SOURCE,
		"postgres", driver)

	if err != nil {
		log.Println("Error while loading migration scripts", err)
	}
	err = m.Drop()
	if err != nil {
		log.Println("Error while running migration down scripts", err)
	}
	err = ExecMigration(connStatement, DB_MIGRATION_SOURCE, true, true)
	if err != nil {
		log.Println("Error while running migration up scripts", err)
	}
}
