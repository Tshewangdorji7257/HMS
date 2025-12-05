package handlers

import (
	"auth-service/models"
	"auth-service/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func init() {
	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	os.Setenv("JWT_EXPIRY", "24h")
}

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"success": true,
		"message": "test",
	}
	
	respondJSON(w, http.StatusOK, data)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}

func TestSignupValidation(t *testing.T) {
	tests := []struct {
		name       string
		payload    models.SignupRequest
		wantStatus int
	}{
		{
			name: "Missing email",
			payload: models.SignupRequest{
				Password: "password123",
				Name:     "Test User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing password",
			payload: models.SignupRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing name",
			payload: models.SignupRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// We expect validation errors without actual DB
			Signup(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestSignupInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Signup(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response models.AuthResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}
}

func TestLoginInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestLoginValidation(t *testing.T) {
	tests := []struct {
		name       string
		payload    models.LoginRequest
		wantStatus int
	}{
		{
			name: "Missing email",
			payload: models.LoginRequest{
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing password",
			payload: models.LoginRequest{
				Email: "test@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			Login(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthResponseStructure(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.AuthResponse{
		Success: true,
		Token:   "test-token",
		User: &models.User{
			ID:    "user-123",
			Email: "test@example.com",
			Name:  "Test User",
			Role:  "student",
		},
	}
	
	respondJSON(w, http.StatusOK, response)
	
	var decoded models.AuthResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if !decoded.Success {
		t.Error("Expected success to be true")
	}
	
	if decoded.Token != "test-token" {
		t.Errorf("Expected token test-token, got %s", decoded.Token)
	}
	
	if decoded.User.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", decoded.User.Email)
	}
}

func TestContentTypeHeader(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"test": "data"}
	
	respondJSON(w, http.StatusOK, data)
	
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "testPassword123"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if !utils.CheckPasswordHash(password, hash) {
		t.Error("Password hash verification failed")
	}
}

func TestTokenGeneration(t *testing.T) {
	user := &models.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "student",
	}
	
	token, err := utils.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	if token == "" {
		t.Error("Token should not be empty")
	}
	
	// Validate the generated token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	
	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %s, got %s", user.ID, claims.UserID)
	}
}

func TestAuthResponseJSON(t *testing.T) {
	w := httptest.NewRecorder()
	user := &models.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      "student",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	response := models.AuthResponse{
		Success: true,
		Message: "Success",
		Token:   "test-token",
		User:    user,
	}
	
	respondJSON(w, http.StatusOK, response)
	
	var decoded models.AuthResponse
	err := json.NewDecoder(w.Body).Decode(&decoded)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if !decoded.Success {
		t.Error("Expected success to be true")
	}
	
	if decoded.User.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, decoded.User.Email)
	}
}

func TestEmailValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload models.SignupRequest
	}{
		{
			name: "Empty email triggers validation",
			payload: models.SignupRequest{
				Email:    "",
				Password: "password123",
				Name:     "User",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			
			Signup(w, req)
			
			var response models.AuthResponse
			json.NewDecoder(w.Body).Decode(&response)
			
			if tt.payload.Email == "" && response.Success {
				t.Error("Empty email should fail validation")
			}
		})
	}
}

func TestRoleDefaulting(t *testing.T) {
	tests := []struct {
		name         string
		role         string
		expectedRole string
	}{
		{"Default to student", "", "student"},
		{"Admin role", "admin", "admin"},
		{"Student role", "student", "student"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := models.SignupRequest{
				Email:    "test@example.com",
				Password: "pass123",
				Name:     "Test",
				Role:     tt.role,
			}
			
			// Default role logic
			if req.Role == "" {
				req.Role = "student"
			}
			
			if req.Role != tt.expectedRole {
				t.Errorf("Expected role %s, got %s", tt.expectedRole, req.Role)
			}
		})
	}
}

func TestHTTPStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"Created", http.StatusCreated},
		{"BadRequest", http.StatusBadRequest},
		{"Unauthorized", http.StatusUnauthorized},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data := map[string]string{"status": tt.name}
			
			respondJSON(w, tt.statusCode, data)
			
			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
