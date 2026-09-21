package config

import (
	"fmt"
	"os"

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

	ProductServiceURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("ORDER_SERVICE_PORT", "8083"),

		DBHost:     os.Getenv("ORDER_DB_HOST"),
		DBPort:     os.Getenv("ORDER_DB_PORT"),
		DBUser:     os.Getenv("ORDER_DB_USER"),
		DBPassword: os.Getenv("ORDER_DB_PASSWORD"),
		DBName:     os.Getenv("ORDER_DB_NAME"),

		ProductServiceURL: getEnv(
			"PRODUCT_SERVICE_URL",
			"http://localhost:8082",
		),
	}

	required := map[string]string{
		"ORDER_DB_HOST":     cfg.DBHost,
		"ORDER_DB_PORT":     cfg.DBPort,
		"ORDER_DB_USER":     cfg.DBUser,
		"ORDER_DB_PASSWORD": cfg.DBPassword,
		"ORDER_DB_NAME":     cfg.DBName,
	}

	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf(
				"required environment variable %s is not set",
				key,
			)
		}
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
