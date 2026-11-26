package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var Files embed.FS

func Up(db *sql.DB) error {
	source, err := iofs.New(Files, ".")
	if err != nil {
		return fmt.Errorf("open embedded migrations: %w", err)
	}

	database, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("initialize PostgreSQL migration driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "postgres", database)
	if err != nil {
		return fmt.Errorf("initialize migration runner: %w", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		sourceErr, databaseErr := migrator.Close()
		return errors.Join(
			fmt.Errorf("apply SQL migrations: %w", err),
			sourceErr,
			databaseErr,
		)
	}

	sourceErr, databaseErr := migrator.Close()
	if err := errors.Join(sourceErr, databaseErr); err != nil {
		return fmt.Errorf("close migration runner: %w", err)
	}

	return nil
}
