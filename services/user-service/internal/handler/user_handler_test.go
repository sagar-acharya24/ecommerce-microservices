package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/auth"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/dto"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/repository"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepository struct {
	users  map[uint]*model.User
	nextID uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[uint]*model.User),
		nextID: 1,
	}
}

func (m *mockUserRepository) Create(user *model.User) error {
	user.ID = m.nextID
	m.nextID++

	m.users[user.ID] = &model.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return nil
}

func (m *mockUserRepository) GetByID(id uint) (*model.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	return &model.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return &model.User{
				ID:           user.ID,
				Name:         user.Name,
				Email:        user.Email,
				PasswordHash: user.PasswordHash,
				CreatedAt:    user.CreatedAt,
				UpdatedAt:    user.UpdatedAt,
			}, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) Update(user *model.User) error {
	existing, exists := m.users[user.ID]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	existing.Name = user.Name
	existing.Email = user.Email
	existing.UpdatedAt = time.Now()

	return nil
}

var _ repository.UserRepository = (*mockUserRepository)(nil)

func newTestUserService(repo repository.UserRepository) service.UserService {
	jwtManager := auth.NewJWTManager(
		"test-secret",
		24*time.Hour,
	)

	return service.NewUserService(repo, jwtManager)
}

func setupTestRouter() (*gin.Engine, *mockUserRepository) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	userService := newTestUserService(repo)
	userHandler := NewUserHandler(userService)

	router := gin.New()

	router.POST("/users/register", userHandler.Register)
	router.POST("/users/login", userHandler.Login)
	router.GET("/users/:id", userHandler.GetUser)
	router.PUT("/users/:id", userHandler.UpdateUser)

	return router, repo
}

func TestUserHandlerRegister(t *testing.T) {
	router, _ := setupTestRouter()

	body := `{
		"name": "Sagar",
		"email": "sagar@example.com",
		"password": "password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/users/register",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			rec.Code,
		)
	}

	var response dto.UserResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID == 0 {
		t.Fatal("expected user ID")
	}

	if response.Email != "sagar@example.com" {
		t.Errorf(
			"expected email sagar@example.com, got %s",
			response.Email,
		)
	}
}

func TestUserHandlerRegisterDuplicateEmail(t *testing.T) {
	router, _ := setupTestRouter()

	body := `{
		"name": "Sagar",
		"email": "sagar@example.com",
		"password": "password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/users/register",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("first registration failed: %d", rec.Code)
	}

	req = httptest.NewRequest(
		http.MethodPost,
		"/users/register",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusConflict,
			rec.Code,
		)
	}
}

func TestUserHandlerLogin(t *testing.T) {
	router, _ := setupTestRouter()

	registerBody := `{
		"name": "Sagar",
		"email": "sagar@example.com",
		"password": "password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/users/register",
		strings.NewReader(registerBody),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("registration failed: %d", rec.Code)
	}

	loginBody := `{
		"email": "sagar@example.com",
		"password": "password123"
	}`

	req = httptest.NewRequest(
		http.MethodPost,
		"/users/login",
		strings.NewReader(loginBody),
	)
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response dto.LoginResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
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
}

func TestUserHandlerLoginInvalidCredentials(t *testing.T) {
	router, _ := setupTestRouter()

	registerBody := `{
		"name": "Sagar",
		"email": "sagar@example.com",
		"password": "password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/users/register",
		strings.NewReader(registerBody),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	loginBody := `{
		"email": "sagar@example.com",
		"password": "wrong-password"
	}`

	req = httptest.NewRequest(
		http.MethodPost,
		"/users/login",
		strings.NewReader(loginBody),
	)
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestUserHandlerGetUser(t *testing.T) {
	router, repo := setupTestRouter()

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: string(passwordHash),
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/1",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response dto.UserResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("expected ID 1, got %d", response.ID)
	}

	if response.Email != "sagar@example.com" {
		t.Errorf(
			"expected email sagar@example.com, got %s",
			response.Email,
		)
	}
}

func TestUserHandlerGetUserNotFound(t *testing.T) {
	router, _ := setupTestRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/999",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestUserHandlerGetUserInvalidID(t *testing.T) {
	router, _ := setupTestRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/abc",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestUserHandlerUpdateUser(t *testing.T) {
	router, repo := setupTestRouter()

	user := &model.User{
		Name:         "Sagar",
		Email:        "sagar@example.com",
		PasswordHash: "hashed-password",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	body := `{
		"name": "Sagar Acharya",
		"email": "sagar.acharya@example.com"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/users/1",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response dto.UserResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Name != "Sagar Acharya" {
		t.Errorf(
			"expected name Sagar Acharya, got %s",
			response.Name,
		)
	}

	if response.Email != "sagar.acharya@example.com" {
		t.Errorf(
			"expected email sagar.acharya@example.com, got %s",
			response.Email,
		)
	}
}

func TestUserHandlerInvalidJSON(t *testing.T) {
	router, _ := setupTestRouter()

	req := httptest.NewRequest(
		http.MethodPost,
		"/users/register",
		strings.NewReader(`invalid-json`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestUserHandlerUpdateDuplicateEmail(t *testing.T) {
	router, repo := setupTestRouter()

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

	body := `{
		"name": "Rahul Updated",
		"email": "sagar@example.com"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/users/2",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusConflict,
			rec.Code,
		)
	}
}
