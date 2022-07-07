package store

import (
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/VladPetriv/tg_scanner/pkg/config"
)

func runMigrations(cfg *config.Config) error {
	if cfg.MigrationsPath == "" {
		return nil
	}

	m, err := migrate.New(cfg.MigrationsPath, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("create migrations error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("up migrations error: %w", err)
	}

	return nil
}
