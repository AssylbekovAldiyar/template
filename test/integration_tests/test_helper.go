package integration_tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgres2 "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Container testcontainers.Container
	DB        *sqlx.DB
}

func setupPostgresContainer(t *testing.T) (testcontainers.Container, string) {
	ctx := context.Background()

	// Define a PostgreSQL container
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:latest"),
		postgres.WithDatabase("example-db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create the container
	connStr, err := container.ConnectionString(ctx, "sslmode=disable", "application_name=tests")
	assert.NoError(t, err)

	return container, connStr
}

func SetupTestDB(t *testing.T) *TestDB {
	// Start the PostgreSQL container
	container, connectionString := setupPostgresContainer(t)

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		t.Fatal(err)
	}

	//run migrations
	if err := runMigrations(db.DB); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return &TestDB{
		Container: container,
		DB:        db,
	}
}

func (tdb *TestDB) Cleanup() {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
	if tdb.Container != nil {
		tdb.Container.Terminate(context.Background())
	}
}

func runMigrations(db *sql.DB) error {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationsPath := fmt.Sprintf("file://%s/../../migration/test", workingDir)

	driver, err := postgres2.WithInstance(db, &postgres2.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run up migrations: %v", err)
	}

	return nil
}

//func getProperPath(path string) string {
//	stack := make([]string, len(path))
//
//}
