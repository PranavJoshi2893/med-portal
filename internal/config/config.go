package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	DBUser          string
	DBName          string
	DBPassword      string
	DBSSLMode       string
	DBHost          string
	DBPort          string
	Pepper          string
	AccessTokenKey  string
	RefreshTokenKey string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %v", err)
	}

	cfg := &Config{
		ServerPort:      os.Getenv("PORT"),
		DBUser:          os.Getenv("POSTGRES_USER"),
		DBName:          os.Getenv("POSTGRES_DB"),
		DBPassword:      os.Getenv("POSTGRES_PASSWORD"),
		DBSSLMode:       os.Getenv("POSTGRES_SSLMODE"),
		DBHost:          os.Getenv("POSTGRES_HOST"),
		DBPort:          os.Getenv("POSTGRES_PORT"),
		Pepper:          os.Getenv("PEPPER"),
		AccessTokenKey:  os.Getenv("ACCESS_TOKEN_KEY"),
		RefreshTokenKey: os.Getenv("REFRESH_TOKEN_KEY"),
	}

	if cfg.Pepper == "" {
		return nil, fmt.Errorf("PEPPER is required")
	}
	if cfg.AccessTokenKey == "" {
		return nil, fmt.Errorf("ACCESS_TOKEN_KEY is required")
	}
	if cfg.RefreshTokenKey == "" {
		return nil, fmt.Errorf("REFRESH_TOKEN_KEY is required")
	}

	return cfg, nil
}
