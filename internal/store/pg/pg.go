package pg

import (
	"database/sql"
	"fmt"

	"github.com/VladPetriv/tg_scanner/pkg/config"
)

type DB struct {
	*sql.DB
}

func Dial(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error while create connection to db: %w", err)
	}

	_, err = db.Exec("SELECT 1;")
	if err != nil {
		return nil, fmt.Errorf("error while send request to db: %w", err)
	}

	return &DB{db}, nil
}
