package service

import (
	"errors"
	"testing"
	"time"

	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/auth"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/dto"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepository struct {
	users         map[uint]*model.User
	nextID        uint
	createErr     error
	updateErr     error
	getByIDErr    error
	getByEmailErr error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[uint]*model.User),
		nextID: 1,
	}
}

func newTestUserService(repo repository.UserRepository) UserService {
	jwtManager := auth.NewJWTManager(
		"test-secret",
		24*time.Hour,
	)

	return NewUserService(repo, jwtManager)
}

func (m *mockUserRepository) Create(user *model.User) error {
	if m.createErr != nil {
		return m.createErr
	}

	user.ID = m.nextID
	m.nextID++

	m.users[user.ID] = &model.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	return nil
}

func (m *mockUserRepository) GetByID(id uint) (*model.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	user, exists := m.users[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	return &model.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}

	for _, user := range m.users {
		if user.Email == email {
			return &model.User{
				ID:           user.ID,
				Name:         user.Name,
				Email:        user.Email,
				PasswordHash: user.PasswordHash,
			}, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) Update(user *model.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}

	existing, exists := m.users[user.ID]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	existing.Name = user.Name
	existing.Email = user.Email
	existing.PasswordHash = user.PasswordHash

	return nil
}

var _ repository.UserRepository = (*mockUserRepository)(nil)

func TestUserServiceRegister(t *testing.T) {
	repo := newMockUserRepository()
	svc := newTestUserService(repo)

	req := dto.RegisterUserRequest{
		Name:     "Sagar",
		Email:    "SAGAR@example.com",
		Password: "password123",
	}

	user, err := svc.Register(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Fatal("expected user ID to be generated")
	}

	if user.Name != "Sagar" {
		t.Errorf("expected name Sagar, got %s", user.Name)
	}

	if user.Email != "sagar@example.com" {
		t.Errorf("expected normalized email sagar@example.com, got %s", user.Email)
	}

	if user.PasswordHash == "" {
		t.Fatal("expected password hash to be generated")
	}

	if user.PasswordHash == req.Password {
		t.Fatal("password must not be stored as plain text")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		t.Fatalf("password hash does not match original password: %v", err)
	}
}

func TestUserServiceRegisterDuplicateEmail(t *testing.T) {
	repo := newMockUserRepository()
	svc := newTestUserService(repo)

	firstRequest := dto.RegisterUserRequest{
		Name:     "Sagar",
		Email:    "sagar@example.com",
		Password: "password123",
	}

	_, err := svc.Register(firstRequest)
	if err != nil {
		t.Fatalf("failed to register first user: %v", err)
	}

	secondRequest := dto.RegisterUserRequest{
		Name:     "Another User",
		Email:    "SAGAR@example.com",
		Password: "password456",
	}

	_, err = svc.Register(secondRequest)
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestUserServiceRegisterValidation(t *testing.T) {
	tests := []struct {
		name string
		req  dto.RegisterUserRequest
	}{
		{
			name: "missing name",
			req: dto.RegisterUserRequest{
				Email:    "sagar@example.com",
				Password: "password123",
			},
		},
		{
			name: "missing email",
			req: dto.RegisterUserRequest{
				Name:     "Sagar",
				Password: "password123",
			},
		},
		{
			name: "missing password",
			req: dto.RegisterUserRequest{
				Name:  "Sagar",
				Email: "sagar@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			svc := newTestUserService(repo)

			_, err := svc.Register(tt.req)
			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestUserServiceGetUserByID(t *testing.T) {
	repo := newMockUserRepository()

	user := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: "hashed-password",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	svc := newTestUserService(repo)

	result, err := svc.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, result.ID)
	}
}

func TestUserServiceGetUserByIDNotFound(t *testing.T) {
	repo := newMockUserRepository()
	svc := newTestUserService(repo)

	_, err := svc.GetUserByID(999)

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserServiceUpdateUser(t *testing.T) {
	repo := newMockUserRepository()

	user := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: "hashed-password",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	svc := newTestUserService(repo)

	req := dto.UpdateUserRequest{
		Name:  "Sagar Acharya",
		Email: "SAGAR.ACHARYA@example.com",
	}

	result, err := svc.UpdateUser(user.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Name != "Sagar Acharya" {
		t.Errorf("expected updated name, got %s", result.Name)
	}

	if result.Email != "sagar.acharya@example.com" {
		t.Errorf("expected normalized email, got %s", result.Email)
	}
}

func TestUserServiceUpdateUserDuplicateEmail(t *testing.T) {
	repo := newMockUserRepository()

	firstUser := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: "hash",
	}

	secondUser := &model.User{
		Name:         "Rahul",
		Email:        "rahul@example.com",
		PasswordHash: "hash",
	}

	if err := repo.Create(firstUser); err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	if err := repo.Create(secondUser); err != nil {
		t.Fatalf("failed to create second user: %v", err)
	}

	svc := newTestUserService(repo)

	req := dto.UpdateUserRequest{
		Name:  "Rahul Updated",
		Email: "sagar@example.com",
	}

	_, err := svc.UpdateUser(secondUser.ID, req)

	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestUserServiceUpdateUserNotFound(t *testing.T) {
	repo := newMockUserRepository()
	svc := newTestUserService(repo)

	req := dto.UpdateUserRequest{
		Name:  "Sagar",
		Email: "sagar@example.com",
	}

	_, err := svc.UpdateUser(999, req)

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserServiceLogin(t *testing.T) {
	repo := newMockUserRepository()

	jwtManager := auth.NewJWTManager(
		"test-secret",
		24*time.Hour,
	)

	svc := NewUserService(repo, jwtManager)

	registerReq := dto.RegisterUserRequest{
		Name:     "Sagar",
		Email:    "sagar@example.com",
		Password: "password123",
	}

	_, err := svc.Register(registerReq)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	loginReq := dto.LoginUserRequest{
		Email:    "SAGAR@example.com",
		Password: "password123",
	}

	response, err := svc.Login(loginReq)
	if err != nil {
		t.Fatalf("expected login to succeed, got %v", err)
	}

	if response.Token == "" {
		t.Fatal("expected JWT token")
	}

	if response.User.Email != "sagar@example.com" {
		t.Errorf(
			"expected email sagar@example.com, got %s",
			response.User.Email,
		)
	}

	if response.User.Name != "Sagar" {
		t.Errorf(
			"expected name Sagar, got %s",
			response.User.Name,
		)
	}

	claims, err := jwtManager.ValidateToken(response.Token)
	if err != nil {
		t.Fatalf("failed to validate generated token: %v", err)
	}

	if claims.UserID == 0 {
		t.Fatal("expected user ID in JWT")
	}
}

func TestUserServiceLoginWrongPassword(t *testing.T) {
	repo := newMockUserRepository()

	jwtManager := auth.NewJWTManager(
		"test-secret",
		24*time.Hour,
	)

	svc := NewUserService(repo, jwtManager)

	_, err := svc.Register(dto.RegisterUserRequest{
		Name:     "Sagar",
		Email:    "sagar@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	_, err = svc.Login(dto.LoginUserRequest{
		Email:    "sagar@example.com",
		Password: "wrong-password",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"expected ErrInvalidCredentials, got %v",
			err,
		)
	}
}

func TestUserServiceLoginUnknownEmail(t *testing.T) {
	repo := newMockUserRepository()

	jwtManager := auth.NewJWTManager(
		"test-secret",
		24*time.Hour,
	)

	svc := NewUserService(repo, jwtManager)

	_, err := svc.Login(dto.LoginUserRequest{
		Email:    "unknown@example.com",
		Password: "password123",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"expected ErrInvalidCredentials, got %v",
			err,
		)
	}
}
