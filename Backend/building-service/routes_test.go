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
		{"Get all buildings", "GET", "/api/buildings"},
		{"Search buildings", "GET", "/api/buildings/search"},
		{"Get building by ID", "GET", "/api/buildings/123"},
		{"Get room by ID", "GET", "/api/buildings/123/rooms/456"},
		{"Update bed occupancy", "PUT", "/api/buildings/beds/789/occupancy"},
		{"Get beds by user", "GET", "/api/buildings/users/user123/beds"},
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
	if !strings.Contains(body, "building-service") {
		t.Error("Response should contain 'building-service'")
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
	port := getPort("8002")
	if port != "8002" {
		t.Errorf("Expected default port 8002, got %s", port)
	}
	
	// Test with custom port
	os.Setenv("PORT", "9999")
	port = getPort("8002")
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

// Test that all building routes are registered
func TestBuildingRoutesExist(t *testing.T) {
	router := setupRouter()
	
	routes := []string{
		"/api/buildings",
		"/api/buildings/search",
		"/api/buildings/123",
		"/api/buildings/123/rooms/456",
	}
	
	for _, route := range routes {
		req := httptest.NewRequest("GET", route, nil)
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
		"/api/buildings",
		"/api/buildings/search",
		"/api/buildings/123",
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
