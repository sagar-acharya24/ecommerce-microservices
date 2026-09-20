package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/config"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/model"
)

type ProductCache interface {
	GetProduct(ctx context.Context, id uint) (*model.Product, error)
	SetProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id uint) error
}

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(cfg *config.Config) *RedisCache {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &RedisCache{
		client: client,
		ttl:    10 * time.Minute,
	}
}

func (r *RedisCache) GetProduct(
	ctx context.Context,
	id uint,
) (*model.Product, error) {
	key := productCacheKey(id)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var product model.Product

	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *RedisCache) SetProduct(
	ctx context.Context,
	product *model.Product,
) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	key := productCacheKey(product.ID)

	return r.client.Set(
		ctx,
		key,
		data,
		r.ttl,
	).Err()
}

func (r *RedisCache) DeleteProduct(
	ctx context.Context,
	id uint,
) error {
	key := productCacheKey(id)

	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func productCacheKey(id uint) string {
	return fmt.Sprintf("product:%d", id)
}
