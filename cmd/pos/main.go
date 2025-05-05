package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"pos-rs/pkg/pos/jsonlog"
	"pos-rs/pkg/pos/model"
	"pos-rs/pkg/pos/vcs"
	"sync"
	"github.com/peterbourgon/ff/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Config struct {
	Port int
	Env  string
	Fill bool
	Migrations string
	DB   struct {
		DSN string
	}
}

var (
	version = vcs.Version()
)

type Application struct {
	Config Config
	Models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	var (
		cfg        Config
		fill       = fs.Bool("fill", false, "Fill database with dummy data")
		migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.Int("port", 8081, "API server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgres://postgres:postgres@localhost:5432/pos_rs?sslmode=disable", "PostgreSQL DSN")
	)

	// Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		logger.PrintFatal(err, nil)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	cfg.Port = *port
	cfg.Env = *env
	cfg.Fill = *fill
	cfg.DB.DSN = *dbDsn
	cfg.Migrations = *migrations

	logger.PrintInfo("starting application with configuration", map[string]string{
		"port":       fmt.Sprintf("%d", cfg.Port),
		"fill":       fmt.Sprintf("%t", cfg.Fill),
		"env":        cfg.Env,
		"db":         cfg.DB.DSN,
		"migrations": cfg.Migrations,
	})

	// Connect to DB
	db, err := OpenDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	// Defer a call to db.Close() so that the connection pool is closed before the main()
	// function exits.
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &Application{
		Config: cfg,
		Models: model.NewModels(db),
		logger: logger,
	}


	// Call app.server() to start the server.
	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func OpenDB(cfg Config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// https://github.com/golang-migrate/migrate?tab=readme-ov-file#use-in-your-go-project
	if cfg.Migrations != "" { 
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}
		m, err := migrate.NewWithDatabaseInstance(
			cfg.Migrations,
			"postgres", driver)
		if err != nil {
			return nil, err
		}
		err = m.Up()
		if err != nil {
			fmt.Fprintf(os.Stderr, "e5: %v\n", err)
			//return nil, err
		}

	}

	return db, nil
}
