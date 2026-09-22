package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port   string

	UserServiceURL    string
	ProductServiceURL string
	OrderServiceURL   string

	JWTSecret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("API_GATEWAY_PORT", "8080"),

		UserServiceURL:    os.Getenv("USER_SERVICE_URL"),
		ProductServiceURL: os.Getenv("PRODUCT_SERVICE_URL"),
		OrderServiceURL:   os.Getenv("ORDER_SERVICE_URL"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	required := map[string]string{
		"USER_SERVICE_URL":    cfg.UserServiceURL,
		"PRODUCT_SERVICE_URL": cfg.ProductServiceURL,
		"ORDER_SERVICE_URL":   cfg.OrderServiceURL,
		"JWT_SECRET":          cfg.JWTSecret,
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
