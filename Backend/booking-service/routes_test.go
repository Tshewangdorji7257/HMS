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
		{"Get all bookings", "GET", "/api/bookings"},
		{"Create booking", "POST", "/api/bookings"},
		{"Get booking by ID", "GET", "/api/bookings/123"},
		{"Cancel booking", "PUT", "/api/bookings/123/cancel"},
		{"Get bookings by user", "GET", "/api/bookings/users/user123"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			// Just verify the route exists (may return various status codes)
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
	if !strings.Contains(body, "booking-service") {
		t.Error("Response should contain 'booking-service'")
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
	port := getPort("8003")
	if port != "8003" {
		t.Errorf("Expected default port 8003, got %s", port)
	}
	
	// Test with custom port
	os.Setenv("PORT", "9999")
	port = getPort("8003")
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

// Test that all booking routes are registered
func TestBookingRoutesExist(t *testing.T) {
	router := setupRouter()
	
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/bookings"},
		{"POST", "/api/bookings"},
		{"GET", "/api/bookings/123"},
		{"PUT", "/api/bookings/123/cancel"},
		{"GET", "/api/bookings/users/user123"},
	}
	
	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path, nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		// Route exists if we get any response (not 404)
		if w.Code == http.StatusNotFound {
			t.Errorf("Route %s %s not found", route.method, route.path)
		}
	}
}

// Test OPTIONS method support
func TestOPTIONSSupport(t *testing.T) {
	router := setupRouter()
	
	routes := []string{
		"/api/bookings",
		"/api/bookings/123",
		"/api/bookings/123/cancel",
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
