package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Phone          string
	Password       string
	DatabaseURL    string
	MigrationsPath string
	ProjectID      string
	StorageBucket  string
	SecretPath     string
	LogLevel       string
}

func Get() (*Config, error) {
	if err := godotenv.Load("configs/.config.env"); err != nil {
		return nil, fmt.Errorf("load env file error: %w", err)
	}

	return &Config{
		Phone:          os.Getenv("PHONE"),
		Password:       os.Getenv("PASSWORD"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		ProjectID:      os.Getenv("PROJECT_ID"),
		StorageBucket:  os.Getenv("STORAGE_BUCKET"),
		SecretPath:     os.Getenv("SECRET_PATH"),
		LogLevel:       os.Getenv("LOG_LEVEL"),
	}, nil
}
