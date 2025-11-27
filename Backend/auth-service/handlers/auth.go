package handlers

import (
	"auth-service/database"
	"auth-service/models"
	"auth-service/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Signup handles user registration
func Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" || req.Name == "" {
		respondJSON(w, http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Error:   "Email, password, and name are required",
		})
		return
	}

	// Default role to student if not provided
	if req.Role == "" {
		req.Role = "student"
	}

	// Check if user already exists
	var existingID string
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
	if err == nil {
		respondJSON(w, http.StatusConflict, models.AuthResponse{
			Success: false,
			Error:   "User with this email already exists",
		})
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking existing user: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.AuthResponse{
			Success: false,
			Error:   "Failed to process password",
		})
		return
	}

	// Create user
	user := &models.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Name:      req.Name,
		Password:  hashedPassword,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = database.DB.Exec(
		"INSERT INTO users (id, email, name, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.ID, user.Email, user.Name, user.Password, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.AuthResponse{
			Success: false,
			Error:   "Failed to create user",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.AuthResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	respondJSON(w, http.StatusCreated, models.AuthResponse{
		Success: true,
		Message: "User created successfully",
		Token:   token,
		User:    user,
	})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		respondJSON(w, http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Error:   "Email and password are required",
		})
		return
	}

	// Get user from database
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, email, name, password, role, created_at, updated_at FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusUnauthorized, models.AuthResponse{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	} else if err != nil {
		log.Printf("Error fetching user: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		respondJSON(w, http.StatusUnauthorized, models.AuthResponse{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.AuthResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	respondJSON(w, http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    &user,
	})
}

// ValidateToken validates JWT token
func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"valid": false,
			"error": "No token provided",
		})
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"valid": false,
			"error": fmt.Sprintf("Invalid token: %v", err),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"valid":  true,
		"claims": claims,
	})
}

// GetUserProfile returns the authenticated user's profile
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "No token provided",
		})
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Invalid token",
		})
		return
	}

	// Get user from database
	var user models.User
	err = database.DB.QueryRow(
		"SELECT id, email, name, role, created_at, updated_at FROM users WHERE id = $1",
		claims.UserID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

// Helper function to send JSON responses
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
