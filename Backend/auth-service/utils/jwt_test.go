package utils

import (
	"auth-service/models"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("JWT_EXPIRY", "24h")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	
	user := &models.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "student",
	}
	
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if token == "" {
		t.Fatal("Expected token to not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("JWT_EXPIRY", "24h")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	
	user := &models.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "student",
	}
	
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %s, got %s", user.ID, claims.UserID)
	}
	
	if claims.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, claims.Email)
	}
	
	if claims.Name != user.Name {
		t.Errorf("Expected Name %s, got %s", user.Name, claims.Name)
	}
	
	if claims.Role != user.Role {
		t.Errorf("Expected Role %s, got %s", user.Role, claims.Role)
	}
}

func TestValidateTokenInvalid(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	
	invalidToken := "invalid.token.string"
	
	_, err := ValidateToken(invalidToken)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	
	user := &models.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "student",
	}
	
	token, _ := GenerateToken(user)
	
	// Change the secret
	jwtSecret = []byte("different-secret")
	
	_, err := ValidateToken(token)
	if err == nil {
		t.Error("Expected error when validating with wrong secret")
	}
}

func TestGenerateTokenCustomExpiry(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("JWT_EXPIRY", "1h")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	
	user := &models.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "admin",
	}
	
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Expected no error validating token, got %v", err)
	}
	
	if claims.Role != "admin" {
		t.Errorf("Expected Role admin, got %s", claims.Role)
	}
}

func TestGenerateTokenInvalidExpiry(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("JWT_EXPIRY", "invalid-duration")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	
	user := &models.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "student",
	}
	
	// Should still generate token with default 24h expiry
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if token == "" {
		t.Fatal("Expected token to be generated with default expiry")
	}
}

func TestValidateTokenExpired(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	
	// Create an expired token manually
	claims := jwt.MapClaims{
		"user_id": "test-user",
		"email":   "test@example.com",
		"name":    "Test User",
		"role":    "student",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)
	
	_, err := ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}
