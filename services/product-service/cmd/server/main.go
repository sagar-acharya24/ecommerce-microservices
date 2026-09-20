package main

import (
	"context"
	"log"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/cache"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/config"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/database"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/handler"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/repository"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/router"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Printf("Starting Product Service on port %s", cfg.Port)

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	log.Println("PostgreSQL connection established")

	if err := db.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("failed to migrate product table: %v", err)
	}

	log.Println("Product database migration completed")

	redisCache := cache.NewRedisCache(cfg)

	if err := redisCache.Ping(context.Background()); err != nil {
		log.Printf("WARNING: Redis connection failed: %v", err)
		log.Println("Product Service will continue without Redis cache")
	} else {
		log.Println("Redis connection established")
	}

	productRepository := repository.NewProductRepository(db)

	productService := service.NewProductService(
		productRepository,
		redisCache,
	)

	productHandler := handler.NewProductHandler(productService)

	r := router.SetupRouter(productHandler)

	log.Println("Product Service initialized successfully")
	log.Printf("Product Service listening on :%s", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start Product Service: %v", err)
	}
}
