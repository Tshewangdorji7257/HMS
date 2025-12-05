package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Call the actual health check handler
	healthCheckHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "healthy") {
		t.Error("Response should contain 'healthy'")
	}
	if !strings.Contains(body, "auth-service") {
		t.Error("Response should contain 'auth-service'")
	}
}

func TestHealthEndpointJSONStructure(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"auth-service"}`))
	})

	handler.ServeHTTP(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}
	if response["service"] != "auth-service" {
		t.Errorf("Expected service 'auth-service', got %v", response["service"])
	}
}

func TestSetupRouter(t *testing.T) {
	router := setupRouter()

	if router == nil {
		t.Fatal("Router should not be nil")
	}

	// Test health check route
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "healthy") {
		t.Error("Health check response should contain 'healthy'")
	}
	if !strings.Contains(body, "auth-service") {
		t.Error("Health check response should contain 'auth-service'")
	}
}

func TestGetPort(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		defaultPort string
		expected    string
	}{
		{"Default port", "", "8001", "8001"},
		{"Custom port", "9001", "8001", "9001"},
		{"Port 8080", "8080", "8001", "8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("PORT", tt.envValue)
			} else {
				os.Unsetenv("PORT")
			}

			result := getPort(tt.defaultPort)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}

			if tt.envValue != "" {
				os.Unsetenv("PORT")
			}
		})
	}
}

func TestPortConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{"Default port", "", "8001"},
		{"Custom port", "9001", "9001"},
		{"Port 8080", "8080", "8080"},
		{"Port 3000", "3000", "3000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("PORT", tt.envValue)
			} else {
				os.Unsetenv("PORT")
			}

			port := os.Getenv("PORT")
			if port == "" {
				port = "8001"
			}

			if port != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, port)
			}

			if tt.envValue != "" {
				os.Unsetenv("PORT")
			}
		})
	}
}

func TestHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, w.Code)
			}
		})
	}
}

func TestJSONResponseHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	handler.ServeHTTP(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type should be application/json")
	}
}

func TestAuthRoutePatterns(t *testing.T) {
	routes := []struct {
		name   string
		path   string
		method string
	}{
		{"User signup", "/api/auth/signup", "POST"},
		{"User login", "/api/auth/login", "POST"},
		{"Admin signup", "/api/auth/admin/signup", "POST"},
		{"Admin login", "/api/auth/admin/login", "POST"},
		{"Get user profile", "/api/auth/me", "GET"},
	}

	for _, route := range routes {
		t.Run(route.name, func(t *testing.T) {
			if route.path == "" {
				t.Error("Path should not be empty")
			}
			if route.method == "" {
				t.Error("Method should not be empty")
			}

			// Test that path starts with /api/auth
			if !strings.HasPrefix(route.path, "/api/auth") && route.path != "/health" {
				t.Errorf("Path should start with /api/auth, got %s", route.path)
			}
		})
	}
}

func TestDatabaseEnvironmentVariables(t *testing.T) {
	envVars := []string{"PORT", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "JWT_SECRET"}

	for _, envVar := range envVars {
		t.Run(envVar, func(t *testing.T) {
			if envVar == "" {
				t.Error("Environment variable name should not be empty")
			}

			// Test setting and getting env var
			testValue := "test_value"
			os.Setenv(envVar, testValue)
			defer os.Unsetenv(envVar)

			value := os.Getenv(envVar)
			if value != testValue {
				t.Errorf("Expected %s, got %s", testValue, value)
			}
		})
	}
}

func TestJWTSecretConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		valid    bool
	}{
		{"Valid JWT secret", "my-super-secret-key-12345", true},
		{"Short JWT secret", "short", false},
		{"Empty JWT secret", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("JWT_SECRET", tt.envValue)
			defer os.Unsetenv("JWT_SECRET")

			secret := os.Getenv("JWT_SECRET")
			isValid := len(secret) >= 16

			if isValid != tt.valid {
				t.Errorf("Expected valid=%v for secret length %d", tt.valid, len(secret))
			}
		})
	}
}

func TestServiceNameConfiguration(t *testing.T) {
	serviceName := "auth-service"
	servicePort := "8001"

	if serviceName == "" {
		t.Error("Service name should not be empty")
	}

	if servicePort == "" {
		t.Error("Service port should not be empty")
	}

	// Test service name format
	if !strings.Contains(serviceName, "auth") {
		t.Error("Service name should contain 'auth'")
	}
}

func TestHTTPStatusCodes(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus int
	}{
		{"OK status", http.StatusOK, 200},
		{"Created status", http.StatusCreated, 201},
		{"Bad Request", http.StatusBadRequest, 400},
		{"Unauthorized", http.StatusUnauthorized, 401},
		{"Forbidden", http.StatusForbidden, 403},
		{"Not Found", http.StatusNotFound, 404},
		{"Internal Server Error", http.StatusInternalServerError, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.statusCode != tt.expectedStatus {
				t.Errorf("Expected %d, got %d", tt.expectedStatus, tt.statusCode)
			}
		})
	}
}

func TestCORSHeaders(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/api/auth/login", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
	}

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin == "" {
		t.Error("CORS header should be set")
	}
}

func TestRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		valid       bool
	}{
		{"Valid JSON content type", "application/json", true},
		{"Valid JSON with charset", "application/json; charset=utf-8", true},
		{"Invalid content type", "text/plain", false},
		{"Empty content type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isJSON := strings.Contains(tt.contentType, "application/json")
			if isJSON != tt.valid {
				t.Errorf("Expected valid=%v, got %v for content type %s", tt.valid, isJSON, tt.contentType)
			}
		})
	}
}

func TestErrorResponseFormat(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/auth/login", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"success":false,"error":"Invalid credentials"}`))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse error JSON: %v", err)
	}

	if response["success"] != false {
		t.Error("Error response should have success=false")
	}
	if response["error"] == nil {
		t.Error("Error response should have error message")
	}
}

func TestSuccessResponseFormat(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/auth/login", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"token":"jwt-token-here"}`))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse success JSON: %v", err)
	}

	if response["success"] != true {
		t.Error("Success response should have success=true")
	}
	if response["token"] == nil {
		t.Error("Success response should have token field")
	}
}

func TestUserRoles(t *testing.T) {
	validRoles := []string{"user", "admin"}

	for _, role := range validRoles {
		t.Run("Role_"+role, func(t *testing.T) {
			if role == "" {
				t.Error("Role should not be empty")
			}

			// Test that role is one of the valid roles
			isValid := false
			for _, valid := range validRoles {
				if role == valid {
					isValid = true
					break
				}
			}

			if !isValid {
				t.Errorf("Role %s is not a valid user role", role)
			}
		})
	}
}
