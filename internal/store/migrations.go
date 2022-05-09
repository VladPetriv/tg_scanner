package store

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func runMigrations(cfg *config.Config) error {
	if cfg.MigrationsPath == "" {
		return nil
	}

	m, err := migrate.New(
		cfg.MigrationsPath,
		fmt.Sprintf(
			"postgresql://%s:%s@%s:5432/%s?sslmode=disable",
			cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgDb,
		),
	)
	if err != nil {
		return fmt.Errorf("create migrations error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("up migrations error: %w", err)
	}

	return nil
}
