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
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", cfg.PgUser, cfg.PgPassword, cfg.PgDb)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("Error while create connection to db: %w", err)
	}
	_, err = db.Exec("SELECT 1;")
	if err != nil {
		return nil, fmt.Errorf("Error while send request to db: %w", err)
	}
	return &DB{db}, nil
}
