package main

import (
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Create router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/auth").Subrouter()

	// Public routes
	api.HandleFunc("/signup", handlers.Signup).Methods("POST", "OPTIONS")
	api.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/validate", handlers.ValidateTokenHandler).Methods("POST", "OPTIONS")

	// Protected routes
	api.HandleFunc("/profile", middleware.AuthMiddleware(handlers.GetUserProfile)).Methods("GET", "OPTIONS")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"auth-service"}`))
	}).Methods("GET")

	// No CORS configuration - API Gateway handles all CORS
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("🚀 Auth Service started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
