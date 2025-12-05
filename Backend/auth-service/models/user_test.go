package models

import (
	"testing"
)

func TestUserModel(t *testing.T) {
	user := User{
		ID:       "user-123",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "hashedpassword",
		Role:     "student",
	}
	
	if user.ID != "user-123" {
		t.Errorf("Expected ID user-123, got %s", user.ID)
	}
	
	if user.Email != "test@example.com" {
		t.Errorf("Expected Email test@example.com, got %s", user.Email)
	}
	
	if user.Role != "student" {
		t.Errorf("Expected Role student, got %s", user.Role)
	}
}

func TestSignupRequest(t *testing.T) {
	req := SignupRequest{
		Email:    "new@example.com",
		Password: "password123",
		Name:     "New User",
		Role:     "admin",
	}
	
	if req.Email != "new@example.com" {
		t.Errorf("Expected Email new@example.com, got %s", req.Email)
	}
	
	if req.Password != "password123" {
		t.Errorf("Expected Password password123, got %s", req.Password)
	}
}

func TestLoginRequest(t *testing.T) {
	req := LoginRequest{
		Email:    "login@example.com",
		Password: "loginpass",
	}
	
	if req.Email != "login@example.com" {
		t.Errorf("Expected Email login@example.com, got %s", req.Email)
	}
}

func TestAuthResponse(t *testing.T) {
	resp := AuthResponse{
		Success: true,
		Token:   "jwt-token",
		User: &User{
			ID:    "user-456",
			Email: "user@example.com",
			Name:  "Response User",
			Role:  "student",
		},
	}
	
	if !resp.Success {
		t.Error("Expected Success to be true")
	}
	
	if resp.Token != "jwt-token" {
		t.Errorf("Expected Token jwt-token, got %s", resp.Token)
	}
	
	if resp.User.ID != "user-456" {
		t.Errorf("Expected User ID user-456, got %s", resp.User.ID)
	}
}

func TestTokenClaims(t *testing.T) {
	claims := TokenClaims{
		UserID: "claim-user",
		Email:  "claim@example.com",
		Name:   "Claim User",
		Role:   "admin",
	}
	
	if claims.UserID != "claim-user" {
		t.Errorf("Expected UserID claim-user, got %s", claims.UserID)
	}
	
	if claims.Role != "admin" {
		t.Errorf("Expected Role admin, got %s", claims.Role)
	}
}
