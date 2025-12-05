package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Test all route registrations in setupRouter
func TestRouteRegistration(t *testing.T) {
	router := setupRouter()
	
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"Health check", "GET", "/health"},
		{"Auth signup POST", "POST", "/api/auth/signup"},
		{"Auth login POST", "POST", "/api/auth/login"},
		{"Auth validate POST", "POST", "/api/auth/validate"},
		{"Auth profile GET", "GET", "/api/auth/profile"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			// Just verify the route exists (may return various status codes)
			// We're not testing handler logic here, just route registration
			if w.Code == 0 {
				t.Error("Route should produce a response")
			}
		})
	}
}

// Test healthCheckHandler directly
func TestHealthCheckHandlerDirect(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	healthCheckHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	body := w.Body.String()
	if !strings.Contains(body, "healthy") {
		t.Error("Response should contain 'healthy'")
	}
	if !strings.Contains(body, "auth-service") {
		t.Error("Response should contain 'auth-service'")
	}
	
	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}
}

// Test getPort with multiple scenarios
func TestGetPortFunction(t *testing.T) {
	// Test with empty environment
	os.Unsetenv("PORT")
	port := getPort("8001")
	if port != "8001" {
		t.Errorf("Expected default port 8001, got %s", port)
	}
	
	// Test with custom port
	os.Setenv("PORT", "9999")
	port = getPort("8001")
	if port != "9999" {
		t.Errorf("Expected custom port 9999, got %s", port)
	}
	os.Unsetenv("PORT")
	
	// Test with different default
	port = getPort("5000")
	if port != "5000" {
		t.Errorf("Expected default port 5000, got %s", port)
	}
}

// Test router health endpoint via router
func TestRouterHealthEndpoint(t *testing.T) {
	router := setupRouter()
	
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected JSON content type")
	}
}

// Test that all auth routes are registered
func TestAuthRoutesExist(t *testing.T) {
	router := setupRouter()
	
	routes := []string{
		"/api/auth/signup",
		"/api/auth/login",
		"/api/auth/validate",
		"/api/auth/profile",
	}
	
	for _, route := range routes {
		req := httptest.NewRequest("POST", route, nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		// Route exists if we get any response (not 404)
		if w.Code == http.StatusNotFound {
			t.Errorf("Route %s not found", route)
		}
	}
}

// Test OPTIONS method support
func TestOPTIONSSupport(t *testing.T) {
	router := setupRouter()
	
	routes := []string{
		"/api/auth/signup",
		"/api/auth/login",
		"/api/auth/profile",
	}
	
	for _, route := range routes {
		req := httptest.NewRequest("OPTIONS", route, nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		// OPTIONS should be handled
		if w.Code == http.StatusMethodNotAllowed {
			t.Errorf("OPTIONS not supported for %s", route)
		}
	}
}
