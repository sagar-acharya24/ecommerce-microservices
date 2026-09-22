package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/auth"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/dto"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	Register(req dto.RegisterUserRequest) (*model.User, error)
	Login(req dto.LoginUserRequest) (*dto.LoginResponse, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*model.User, error)
}

type userService struct {
	repo       repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(
	repo repository.UserRepository,
	jwtManager *auth.JWTManager,
) UserService {
	return &userService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *userService) Register(req dto.RegisterUserRequest) (*model.User, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	_, err := s.repo.GetByEmail(email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) Login(
	req dto.LoginUserRequest,
) (*dto.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" {
		return nil, errors.New("email is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) UpdateUser(
	id uint,
	req dto.UpdateUserRequest,
) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	existingUser, err := s.repo.GetByEmail(email)
	if err == nil && existingUser.ID != id {
		return nil, ErrEmailAlreadyExists
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	user.Name = name
	user.Email = email

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
