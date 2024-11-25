package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Scorzoner/effective-mobile-test/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Open(config config.Config) (db *sql.DB, err error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return
	}

	return
}

func RunMigrations(db *sql.DB) error {
	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations", "postgres", migrationDriver)
	if err != nil {
		return err
	}

	err = migrator.Up()
	return err
}
