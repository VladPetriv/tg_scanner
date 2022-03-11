package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Phone    string
	Password string
	Limit    int
}

func Get() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_LOAD_ENV_FILE:%w", err)
	}

	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))

	return &Config{
		Phone:    os.Getenv("PHONE"),
		Password: os.Getenv("PASSWORD"),
		Limit:    limit,
	}, nil
}
