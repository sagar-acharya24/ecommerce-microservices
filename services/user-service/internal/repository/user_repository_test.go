package repository

import (
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/model"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		t.Name(),
	)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestUserRepositoryCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: "hashed-password",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Fatal("expected user ID to be generated")
	}

	if user.Name != "Sagar" {
		t.Errorf("expected name Sagar, got %s", user.Name)
	}
}

func TestUserRepositoryGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: "hashed-password",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	result, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to get user by ID: %v", err)
	}

	if result.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, result.ID)
	}

	if result.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, result.Email)
	}
}

func TestUserRepositoryGetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: "hashed-password",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	result, err := repo.GetByEmail("sagar@example.com")
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}

	if result.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, result.ID)
	}
}

func TestUserRepositoryUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: "hashed-password",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user.Name = "Sagar Acharya"

	if err := repo.Update(user); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	result, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to get updated user: %v", err)
	}

	if result.Name != "Sagar Acharya" {
		t.Errorf("expected name Sagar Acharya, got %s", result.Name)
	}
}

func TestUserRepositoryGetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.GetByID(9999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestUserRepositoryGetByEmailNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.GetByEmail("does-not-exist@example.com")
	if err == nil {
		t.Fatal("expected error for non-existent email")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}
