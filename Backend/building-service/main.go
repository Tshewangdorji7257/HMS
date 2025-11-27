package main

import (
	"building-service/database"
	"building-service/handlers"
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
	api := router.PathPrefix("/api/buildings").Subrouter()

	// Building routes
	api.HandleFunc("", handlers.GetAllBuildings).Methods("GET", "OPTIONS")
	api.HandleFunc("/search", handlers.SearchBuildings).Methods("GET", "OPTIONS")
	api.HandleFunc("/{id}", handlers.GetBuildingByID).Methods("GET", "OPTIONS")
	api.HandleFunc("/{id}/rooms/{roomId}", handlers.GetRoomByID).Methods("GET", "OPTIONS")
	api.HandleFunc("/beds/{bedId}/occupancy", handlers.UpdateBedOccupancy).Methods("PUT", "OPTIONS")
	api.HandleFunc("/users/{userId}/beds", handlers.GetBedsByUserID).Methods("GET", "OPTIONS")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"building-service"}`))
	}).Methods("GET")

	// No CORS configuration - API Gateway handles all CORS
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("🚀 Building Service started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
