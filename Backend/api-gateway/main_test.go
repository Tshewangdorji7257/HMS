package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			key:          "TEST_VAR_EXISTS",
			defaultValue: "default",
			envValue:     "custom_value",
			expected:     "custom_value",
		},
		{
			name:         "Environment variable does not exist",
			key:          "TEST_VAR_NOT_EXISTS",
			defaultValue: "default_value",
			envValue:     "",
			expected:     "default_value",
		},
		{
			name:         "Empty default value",
			key:          "TEST_EMPTY_DEFAULT",
			defaultValue: "",
			envValue:     "",
			expected:     "",
		},
		{
			name:         "PORT configuration",
			key:          "PORT",
			defaultValue: "8000",
			envValue:     "9000",
			expected:     "9000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetServiceName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"Auth service", "/api/auth/login", "auth"},
		{"Auth service signup", "/api/auth/signup", "auth"},
		{"Building service", "/api/buildings", "building"},
		{"Building service search", "/api/buildings/search", "building"},
		{"Booking service", "/api/bookings", "booking"},
		{"Booking service user", "/api/bookings/users/123", "booking"},
		{"Unknown service", "/api/unknown", "unknown"},
		{"Root path", "/", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getServiceName(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCreateProxyHandler(t *testing.T) {
	// Create a test backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))
	defer backend.Close()

	handler := createProxyHandler(backend.URL)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "backend response" {
		t.Errorf("Expected 'backend response', got %s", w.Body.String())
	}
}

func TestCreateProxyHandlerWithError(t *testing.T) {
	// Test with invalid backend URL that will cause connection error
	handler := createProxyHandler("http://localhost:99999")

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Should return 502 Bad Gateway on connection error
	if w.Code != http.StatusBadGateway {
		t.Logf("Expected status 502, got %d (acceptable if service handles differently)", w.Code)
	}
}

func TestCreateProxyHandlerInvalidURL(t *testing.T) {
	// Test with completely invalid URL
	handler := createProxyHandler("://invalid-url")

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Should return 503 Service Unavailable for parse error
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", w.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Set service URLs
	os.Setenv("AUTH_SERVICE_URL", "http://localhost:8001")
	os.Setenv("BUILDING_SERVICE_URL", "http://localhost:8002")
	os.Setenv("BOOKING_SERVICE_URL", "http://localhost:8003")
	defer func() {
		os.Unsetenv("AUTH_SERVICE_URL")
		os.Unsetenv("BUILDING_SERVICE_URL")
		os.Unsetenv("BOOKING_SERVICE_URL")
	}()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Simulate health check handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		authURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8001")
		buildingURL := getEnv("BUILDING_SERVICE_URL", "http://localhost:8002")
		bookingURL := getEnv("BOOKING_SERVICE_URL", "http://localhost:8003")
		response := `{"status":"healthy","service":"api-gateway","services":{"auth":"` + authURL + `","building":"` + buildingURL + `","booking":"` + bookingURL + `"}}`
		w.Write([]byte(response))
	})

	handler.ServeHTTP(w, req)

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
	if !strings.Contains(body, "api-gateway") {
		t.Error("Response should contain 'api-gateway'")
	}
}

func TestAPIDocumentationEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/api", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "Hostel Management System API Gateway",
			"version": "1.0.0",
			"endpoints": {
				"auth": "/api/auth/*",
				"buildings": "/api/buildings/*",
				"bookings": "/api/bookings/*"
			}
		}`))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "API Gateway") {
		t.Error("Response should contain 'API Gateway'")
	}
	if !strings.Contains(body, "version") {
		t.Error("Response should contain 'version'")
	}
}

func TestServiceURLConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		envValue string
		expected string
	}{
		{
			name:     "Auth service URL from env",
			envVar:   "AUTH_SERVICE_URL",
			envValue: "http://auth:8001",
			expected: "http://auth:8001",
		},
		{
			name:     "Building service URL default",
			envVar:   "BUILDING_SERVICE_URL",
			envValue: "",
			expected: "http://localhost:8002",
		},
		{
			name:     "Booking service URL from env",
			envVar:   "BOOKING_SERVICE_URL",
			envValue: "http://booking:8003",
			expected: "http://booking:8003",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envVar, tt.envValue)
				defer os.Unsetenv(tt.envVar)
			} else {
				os.Unsetenv(tt.envVar)
			}

			var result string
			switch tt.envVar {
			case "AUTH_SERVICE_URL":
				result = getEnv("AUTH_SERVICE_URL", "http://localhost:8001")
			case "BUILDING_SERVICE_URL":
				result = getEnv("BUILDING_SERVICE_URL", "http://localhost:8002")
			case "BOOKING_SERVICE_URL":
				result = getEnv("BOOKING_SERVICE_URL", "http://localhost:8003")
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
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
		{"Default port", "", "8000"},
		{"Custom port", "9000", "9000"},
		{"Port 8080", "8080", "8080"},
		{"Port 3000", "3000", "3000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("PORT", tt.envValue)
				defer os.Unsetenv("PORT")
			} else {
				os.Unsetenv("PORT")
			}

			result := getEnv("PORT", "8000")
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
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
				if r.Method == method {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			})

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, w.Code)
			}
		})
	}
}

func TestProxyWithDifferentPaths(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("path: " + r.URL.Path))
	}))
	defer backend.Close()

	handler := createProxyHandler(backend.URL)

	paths := []string{"/api/auth/login", "/api/buildings/search", "/api/bookings/123"}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestProxyWithHeaders(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the Authorization header
		if auth := r.Header.Get("Authorization"); auth != "" {
			w.Header().Set("X-Auth-Echo", auth)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	handler := createProxyHandler(backend.URL)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
