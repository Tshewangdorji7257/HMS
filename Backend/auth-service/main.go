package main

import (
	"auth-service/consul"
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	// Initialize Consul
	if err := consul.InitConsul(); err != nil {
		log.Printf("⚠️  Failed to initialize Consul: %v", err)
	} else {
		// Register service with Consul
		if err := consul.RegisterService(); err != nil {
			log.Printf("⚠️  Failed to register service with Consul: %v", err)
		}
		// Deregister on shutdown
		defer consul.DeregisterService()
	}

	// Create router
	router := setupRouter()

	// No CORS configuration - API Gateway handles all CORS
	// Start server
	port := getPort("8001")

	// Setup graceful shutdown
	go func() {
		log.Printf("🚀 Auth Service started on port %s", port)
		if err := http.ListenAndServe(":"+port, router); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Auth Service...")
}

// setupRouter configures and returns the HTTP router with all routes
func setupRouter() *mux.Router {
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
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	return router
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"auth-service"}`))
}

// getPort returns the port from environment or default
func getPort(defaultPort string) string {
	port := os.Getenv("PORT")
	if port == "" {
		return defaultPort
	}
	return port
}
