package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Phone         string
	Password      string
	PgUser        string
	PgPassword    string
	PgDb          string
	ProjectID     string
	StorageBucket string
	SecretPath    string
}

func Get() (*Config, error) {
	if err := godotenv.Load("configs/.config.env"); err != nil {
		return nil, fmt.Errorf("load env file error: %w", err)
	}

	return &Config{
		Phone:         os.Getenv("PHONE"),
		Password:      os.Getenv("PASSWORD"),
		PgUser:        os.Getenv("POSTGRES_USER"),
		PgPassword:    os.Getenv("POSTGRES_PASSWORD"),
		PgDb:          os.Getenv("POSTGRES_DB"),
		ProjectID:     os.Getenv("PROJECT_ID"),
		StorageBucket: os.Getenv("STORAGE_BUCKET"),
		SecretPath:    os.Getenv("SECRET_PATH"),
	}, nil
}
