package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Phone         string
	Password      string
	ProjectID     string
	StorageBucket string
	SecretPath    string
	LogLevel      string
	LogFilename   string
	RedisPassword string
	RedisAddr     string
	KafkaAddr     string
}

func Get() (*Config, error) {
	if err := godotenv.Load("configs/.config.env"); err != nil {
		return nil, fmt.Errorf("load env file error: %w", err)
	}

	return &Config{
		Phone:         os.Getenv("PHONE"),
		Password:      os.Getenv("PASSWORD"),
		ProjectID:     os.Getenv("PROJECT_ID"),
		StorageBucket: os.Getenv("STORAGE_BUCKET"),
		SecretPath:    os.Getenv("SECRET_PATH"),
		LogLevel:      os.Getenv("LOG_LEVEL"),
		LogFilename:   os.Getenv("LOG_FILENAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		KafkaAddr:     os.Getenv("KAFKA_ADDR"),
	}, nil
}
