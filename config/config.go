package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Phone    string
	Password string
}

func GetConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("ERROR_WHILE_LOAD_ENV_FILE")
	}

	return &Config{
		Phone:    os.Getenv("PHONE"),
		Password: os.Getenv("PASSWORD"),
	}, nil
}
