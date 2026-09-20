package cache

import (
	"context"
	"testing"

	"github.com/joho/godotenv"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/config"
)

func TestRedisConnection(t *testing.T) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		t.Fatalf("failed to load .env: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	redisCache := NewRedisCache(cfg)

	ctx := context.Background()

	if err := redisCache.Ping(ctx); err != nil {
		t.Fatalf("failed to connect to Redis: %v", err)
	}

	t.Log("Redis connection successful")
}
