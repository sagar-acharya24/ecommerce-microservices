package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port   string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func Load() (*Config, error) {
	// Load .env for local development.
	// In Docker/production, environment variables can be
	// provided directly by the environment.
	_ = godotenv.Load()

	redisDB := 0

	if value := os.Getenv("REDIS_DB"); value != "" {
		parsedDB, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
		}

		redisDB = parsedDB
	}

	cfg := &Config{
		AppEnv: os.Getenv("APP_ENV"),
		Port:   os.Getenv("PRODUCT_SERVICE_PORT"),

		DBHost:     os.Getenv("PRODUCT_DB_HOST"),
		DBPort:     os.Getenv("PRODUCT_DB_PORT"),
		DBUser:     os.Getenv("PRODUCT_DB_USER"),
		DBPassword: os.Getenv("PRODUCT_DB_PASSWORD"),
		DBName:     os.Getenv("PRODUCT_DB_NAME"),

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
	}

	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}

	if cfg.Port == "" {
		cfg.Port = "8082"
	}

	if cfg.DBHost == "" {
		return nil, fmt.Errorf("PRODUCT_DB_HOST is required")
	}

	if cfg.DBPort == "" {
		return nil, fmt.Errorf("PRODUCT_DB_PORT is required")
	}

	if cfg.DBUser == "" {
		return nil, fmt.Errorf("PRODUCT_DB_USER is required")
	}

	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("PRODUCT_DB_PASSWORD is required")
	}

	if cfg.DBName == "" {
		return nil, fmt.Errorf("PRODUCT_DB_NAME is required")
	}

	return cfg, nil
}
