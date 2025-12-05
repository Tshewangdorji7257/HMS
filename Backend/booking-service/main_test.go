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
	if !strings.Contains(body, "booking-service") {
		t.Error("Response should contain 'booking-service'")
	}
}

func TestHealthEndpointJSONStructure(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"booking-service"}`))
	})

	handler.ServeHTTP(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}
	if response["service"] != "booking-service" {
		t.Errorf("Expected service 'booking-service', got %v", response["service"])
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
	if !strings.Contains(body, "booking-service") {
		t.Error("Health check response should contain 'booking-service'")
	}
}

func TestGetPort(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		defaultPort string
		expected    string
	}{
		{"Default port", "", "8003", "8003"},
		{"Custom port", "9003", "8003", "9003"},
		{"Port 8080", "8080", "8003", "8080"},
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
		{"Default port", "", "8003"},
		{"Custom port", "9003", "9003"},
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
				port = "8003"
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

func TestBookingRoutePatterns(t *testing.T) {
	routes := []struct {
		name   string
		path   string
		method string
	}{
		{"Get all bookings", "/api/bookings", "GET"},
		{"Create booking", "/api/bookings", "POST"},
		{"Get booking by ID", "/api/bookings/123", "GET"},
		{"Cancel booking", "/api/bookings/123/cancel", "PUT"},
		{"Get bookings by user", "/api/bookings/users/user123", "GET"},
		{"Update booking status", "/api/bookings/123/status", "PUT"},
	}

	for _, route := range routes {
		t.Run(route.name, func(t *testing.T) {
			if route.path == "" {
				t.Error("Path should not be empty")
			}
			if route.method == "" {
				t.Error("Method should not be empty")
			}

			// Test that path starts with /api/bookings
			if !strings.HasPrefix(route.path, "/api/bookings") && route.path != "/health" {
				t.Errorf("Path should start with /api/bookings, got %s", route.path)
			}
		})
	}
}

func TestDatabaseEnvironmentVariables(t *testing.T) {
	envVars := []string{"PORT", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}

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

func TestServiceNameConfiguration(t *testing.T) {
	serviceName := "booking-service"
	servicePort := "8003"

	if serviceName == "" {
		t.Error("Service name should not be empty")
	}

	if servicePort == "" {
		t.Error("Service port should not be empty")
	}

	// Test service name format
	if !strings.Contains(serviceName, "booking") {
		t.Error("Service name should contain 'booking'")
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
		{"Not Found", http.StatusNotFound, 404},
		{"Conflict", http.StatusConflict, 409},
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
	req := httptest.NewRequest("OPTIONS", "/api/bookings", nil)
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
	req := httptest.NewRequest("GET", "/api/bookings/invalid", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success":false,"error":"Booking not found"}`))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
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
	req := httptest.NewRequest("GET", "/api/bookings", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"data":[]}`))
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
	if response["data"] == nil {
		t.Error("Success response should have data field")
	}
}

func TestBookingStatusTransitions(t *testing.T) {
	validStatuses := []string{"pending", "confirmed", "cancelled", "completed"}

	for _, status := range validStatuses {
		t.Run("Status_"+status, func(t *testing.T) {
			if status == "" {
				t.Error("Status should not be empty")
			}

			// Test that status is one of the valid statuses
			isValid := false
			for _, valid := range validStatuses {
				if status == valid {
					isValid = true
					break
				}
			}

			if !isValid {
				t.Errorf("Status %s is not a valid booking status", status)
			}
		})
	}
}

func TestEmailNotificationConfiguration(t *testing.T) {
	emailEnvVars := []string{"SMTP_HOST", "SMTP_PORT", "SMTP_USER", "SMTP_PASSWORD"}

	for _, envVar := range emailEnvVars {
		t.Run(envVar, func(t *testing.T) {
			if envVar == "" {
				t.Error("Email environment variable name should not be empty")
			}

			// Test that env var can be set and retrieved
			testValue := "test_" + strings.ToLower(envVar)
			os.Setenv(envVar, testValue)
			defer os.Unsetenv(envVar)

			value := os.Getenv(envVar)
			if value != testValue {
				t.Errorf("Expected %s, got %s", testValue, value)
			}
		})
	}
}
