package main

import (
	"booking-service/database"
	"booking-service/handlers"
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
	api := router.PathPrefix("/api/bookings").Subrouter()

	// Booking routes
	api.HandleFunc("", handlers.GetAllBookings).Methods("GET", "OPTIONS")
	api.HandleFunc("", handlers.CreateBooking).Methods("POST", "OPTIONS")
	api.HandleFunc("/{id}", handlers.GetBookingByID).Methods("GET", "OPTIONS")
	api.HandleFunc("/{id}/cancel", handlers.CancelBooking).Methods("PUT", "OPTIONS")
	api.HandleFunc("/users/{userId}", handlers.GetBookingsByUserID).Methods("GET", "OPTIONS")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"booking-service"}`))
	}).Methods("GET")

	// No CORS configuration - API Gateway handles all CORS
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	log.Printf("🚀 Booking Service started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
