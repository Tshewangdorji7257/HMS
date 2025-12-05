package middleware

import (
	"auth-service/models"
	"auth-service/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	
	// Create a test user and token
	user := &models.User{
		ID:    "test-user",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "student",
	}
	
	token, _ := utils.GenerateToken(user)
	
	// Test with valid token
	handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestAuthMiddlewareNoToken(t *testing.T) {
	handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddlewareTokenWithoutBearer(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	
	user := &models.User{
		ID:    "test-user",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "student",
	}
	
	token, _ := utils.GenerateToken(user)
	
	handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestRequireRole(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	
	// Test with admin role
	adminUser := &models.User{
		ID:    "admin-user",
		Email: "admin@example.com",
		Name:  "Admin User",
		Role:  "admin",
	}
	
	adminToken, _ := utils.GenerateToken(adminUser)
	
	handler := RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Admin Access"))
	})
	
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 for admin, got %d", rr.Code)
	}
}

func TestRequireRoleInsufficientPermissions(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	
	// Test with student trying to access admin route
	studentUser := &models.User{
		ID:    "student-user",
		Email: "student@example.com",
		Name:  "Student User",
		Role:  "student",
	}
	
	studentToken, _ := utils.GenerateToken(studentUser)
	
	handler := RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+studentToken)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

func TestRequireRoleNoToken(t *testing.T) {
	handler := RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	req := httptest.NewRequest("GET", "/admin", nil)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestRequireRoleInvalidToken(t *testing.T) {
	handler := RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}
