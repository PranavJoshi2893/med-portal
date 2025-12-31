package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
}

func Load() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %v", err)
	}

	return &Config{
		ServerPort: os.Getenv("PORT"),
	}, nil
}
