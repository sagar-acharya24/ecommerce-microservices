package database

import (
	"fmt"

	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to PostgreSQL: %w",
			err,
		)
	}

	return db, nil
}
