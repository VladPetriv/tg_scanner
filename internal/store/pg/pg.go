package pg

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func Dial(cfg *config.Config) (*DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error while create connection to db: %w", err)
	}

	_, err = db.Exec("SELECT 1;")
	if err != nil {
		return nil, fmt.Errorf("error while send request to db: %w", err)
	}

	return &DB{db}, nil
}
