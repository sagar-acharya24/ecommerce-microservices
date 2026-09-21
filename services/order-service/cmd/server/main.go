package main

import (
	"log"

	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/client"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/config"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/database"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/handler"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/repository"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/router"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Println("Starting Order Service")

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	log.Println("PostgreSQL connection established")

	if err := db.AutoMigrate(&model.Order{}); err != nil {
		log.Fatalf("failed to migrate order table: %v", err)
	}

	log.Println("Order database migration completed")

	orderRepository := repository.NewOrderRepository(db)

	productClient := client.NewProductClient(
		cfg.ProductServiceURL,
	)

	orderService := service.NewOrderService(
		orderRepository,
		productClient,
	)

	orderHandler := handler.NewOrderHandler(orderService)

	r := router.SetupRouter(orderHandler)

	log.Println("Order Service initialized successfully")
	log.Printf("Order Service listening on :%s", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start Order Service: %v", err)
	}
}
