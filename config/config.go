package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Phone       string
	Password    string
	Limit       int
	PG_USER     string
	PG_PASSWORD string
	PG_DB       string
}

func Get() (*Config, error) {
	if err := godotenv.Load("configs/.config.env"); err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_LOAD_ENV_FILE:%w", err)
	}

	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))

	return &Config{
		Phone:       os.Getenv("PHONE"),
		Password:    os.Getenv("PASSWORD"),
		Limit:       limit,
		PG_USER:     os.Getenv("POSTGRES_USER"),
		PG_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		PG_DB:       os.Getenv("POSTGRES_DB"),
	}, nil
}
