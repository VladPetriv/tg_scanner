package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Phone      string
	Password   string
	Limit      int
	PgUser     string
	PgPassword string
	PgDb       string
}

func Get() (*Config, error) {
	if err := godotenv.Load("configs/.config.env"); err != nil {
		return nil, fmt.Errorf("load env file error: %w", err)
	}

	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))

	return &Config{
		Phone:      os.Getenv("PHONE"),
		Password:   os.Getenv("PASSWORD"),
		Limit:      limit,
		PgUser:     os.Getenv("POSTGRES_USER"),
		PgPassword: os.Getenv("POSTGRES_PASSWORD"),
		PgDb:       os.Getenv("POSTGRES_DB"),
	}, nil
}
