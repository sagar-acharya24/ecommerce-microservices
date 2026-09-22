package main

import (
	"log"
	"time"

	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/auth"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/config"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/database"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/handler"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/repository"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/router"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/service"
)

func main() {
	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Println("Starting User Service")
	log.Printf("Environment: %s", cfg.AppEnv)
	log.Printf("Port: %s", cfg.Port)
	log.Printf("Database: %s", cfg.DBName)

	// Connect to PostgreSQL.
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	log.Println("PostgreSQL connection established")

	// Run database migrations.
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate user database: %v", err)
	}

	log.Println("User database migration completed")

	// Create repository.
	userRepository := repository.NewUserRepository(db)

	// Create JWT manager.
	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
		24*time.Hour,
	)

	// Create service.
	userService := service.NewUserService(
		userRepository,
		jwtManager,
	)

	// Create handler.
	userHandler := handler.NewUserHandler(userService)

	// Create router.
	r := router.SetupRouter(userHandler, jwtManager)

	log.Printf("User Service initialized successfully")
	log.Printf("User Service listening on :%s", cfg.Port)

	// Start HTTP server.
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start User Service: %v", err)
	}
}
