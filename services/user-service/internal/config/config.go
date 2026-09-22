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

	JWTSecret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("USER_SERVICE_PORT", "8081"),

		DBHost:     os.Getenv("USER_DB_HOST"),
		DBPort:     os.Getenv("USER_DB_PORT"),
		DBUser:     os.Getenv("USER_DB_USER"),
		DBPassword: os.Getenv("USER_DB_PASSWORD"),
		DBName:     os.Getenv("USER_DB_NAME"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	required := map[string]string{
		"USER_DB_HOST":     cfg.DBHost,
		"USER_DB_PORT":     cfg.DBPort,
		"USER_DB_USER":     cfg.DBUser,
		"USER_DB_PASSWORD": cfg.DBPassword,
		"USER_DB_NAME":     cfg.DBName,
		"JWT_SECRET":       cfg.JWTSecret,
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
