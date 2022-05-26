package pg

import (
	"database/sql"
	"fmt"

	"github.com/VladPetriv/tg_scanner/config"
)

type DB struct {
	*sql.DB
}

func Dial(cfg config.Config) (*DB, error) {
	var connectionString string

	if cfg.DatabaseURL == "" {
		connectionString = fmt.Sprintf(
			"user=%s password=%s host=%s dbname=%s sslmode=disable",
			cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgDb,
		)
	} else {
		connectionString = cfg.DatabaseURL
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error while create connection to db: %w", err)
	}

	_, err = db.Exec("SELECT 1;")
	if err != nil {
		return nil, fmt.Errorf("error while send request to db: %w", err)
	}

	return &DB{db}, nil
}
