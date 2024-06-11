package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("failed to load .env file: %v\n", err)
		os.Exit(1)
	}
}

type Config struct {
	LogLevel     slog.Level `env:"LOG_LEVEL,required"`
	DBConnString string     `env:"DB_CONN_STRING,required"`
	ServerAddr   string     `env:"SERVER_ADDR,required"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("env: %w", err)
	}

	return cfg, nil
}
