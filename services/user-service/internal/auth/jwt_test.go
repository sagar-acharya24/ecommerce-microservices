package auth

import (
	"testing"
	"time"
)

func TestJWTManagerGenerateAndValidateToken(t *testing.T) {
	manager := NewJWTManager(
		"test-secret",
		24*time.Hour,
	)

	token, err := manager.GenerateToken(123)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("expected token to be generated")
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != 123 {
		t.Errorf(
			"expected user ID 123, got %d",
			claims.UserID,
		)
	}
}

func TestJWTManagerInvalidToken(t *testing.T) {
	manager := NewJWTManager(
		"test-secret",
		24*time.Hour,
	)

	_, err := manager.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("expected invalid token error")
	}
}

func TestJWTManagerWrongSecret(t *testing.T) {
	manager := NewJWTManager(
		"test-secret",
		24*time.Hour,
	)

	token, err := manager.GenerateToken(123)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	anotherManager := NewJWTManager(
		"wrong-secret",
		24*time.Hour,
	)

	_, err = anotherManager.ValidateToken(token)
	if err == nil {
		t.Fatal("expected token validation to fail with wrong secret")
	}
}

func TestJWTManagerExpiredToken(t *testing.T) {
	manager := NewJWTManager(
		"test-secret",
		-1*time.Hour,
	)

	token, err := manager.GenerateToken(123)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = manager.ValidateToken(token)
	if err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}
