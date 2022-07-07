package store

import (
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/VladPetriv/tg_scanner/pkg/config"
)

func runMigrations(cfg *config.Config) error {
	var connectionString string

	if cfg.MigrationsPath == "" {
		return nil
	}

	if cfg.DatabaseURL == "" {
		connectionString = fmt.Sprintf(
			"postgresql://%s:%s@%s:5432/%s?sslmode=disable",
			cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgDb,
		)
	} else {
		connectionString = cfg.DatabaseURL
	}

	m, err := migrate.New(cfg.MigrationsPath, connectionString)
	if err != nil {
		return fmt.Errorf("create migrations error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("up migrations error: %w", err)
	}

	return nil
}
