package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Phone    string
	Password string
	Limit    int
}

func GetConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("ERROR_WHILE_LOAD_ENV_FILE")
	}
	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))

	return &Config{
		Phone:    os.Getenv("PHONE"),
		Password: os.Getenv("PASSWORD"),
		Limit:    limit,
	}, nil
}
