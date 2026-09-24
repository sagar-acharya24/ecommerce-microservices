package main

import (
	"log"

	"github.com/sagar-acharya24/ecommerce-microservices/api-gateway/internal/config"
	"github.com/sagar-acharya24/ecommerce-microservices/api-gateway/internal/proxy"
	"github.com/sagar-acharya24/ecommerce-microservices/api-gateway/internal/router"
)

func main() {
	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Println("Starting API Gateway")
	log.Printf("Environment: %s", cfg.AppEnv)
	log.Printf("Port: %s", cfg.Port)

	// Create User Service proxy.
	userProxy, err := proxy.NewServiceProxy(cfg.UserServiceURL)
	if err != nil {
		log.Fatalf("failed to create user service proxy: %v", err)
	}

	// Create Product Service proxy.
	productProxy, err := proxy.NewServiceProxy(cfg.ProductServiceURL)
	if err != nil {
		log.Fatalf("failed to create product service proxy: %v", err)
	}

	// Create Order Service proxy.
	orderProxy, err := proxy.NewServiceProxy(cfg.OrderServiceURL)
	if err != nil {
		log.Fatalf("failed to create order service proxy: %v", err)
	}

	// Create router.
	r := router.SetupRouter(
		userProxy,
		productProxy,
		orderProxy,
		cfg.JWTSecret,
	)

	log.Println("API Gateway initialized successfully")
	log.Printf("API Gateway listening on :%s", cfg.Port)

	// Start HTTP server.
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start API Gateway: %v", err)
	}
}
