package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var (
	authServiceURL     string
	buildingServiceURL string
	bookingServiceURL  string
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load service URLs
	authServiceURL = getEnv("AUTH_SERVICE_URL", "http://localhost:8001")
	buildingServiceURL = getEnv("BUILDING_SERVICE_URL", "http://localhost:8002")
	bookingServiceURL = getEnv("BOOKING_SERVICE_URL", "http://localhost:8003")

	// Create router
	router := mux.NewRouter()

	// Auth service routes
	router.PathPrefix("/api/auth").HandlerFunc(createProxyHandler(authServiceURL))

	// Building service routes
	router.PathPrefix("/api/buildings").HandlerFunc(createProxyHandler(buildingServiceURL))

	// Booking service routes
	router.PathPrefix("/api/bookings").HandlerFunc(createProxyHandler(bookingServiceURL))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"api-gateway","services":{"auth":"` + authServiceURL + `","building":"` + buildingServiceURL + `","booking":"` + bookingServiceURL + `"}}`))
	}).Methods("GET")

	// API documentation
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "Hostel Management System API Gateway",
			"version": "1.0.0",
			"endpoints": {
				"auth": "/api/auth/*",
				"buildings": "/api/buildings/*",
				"bookings": "/api/bookings/*"
			},
			"documentation": {
				"auth_service": "Authentication and user management",
				"building_service": "Building, room, and bed management",
				"booking_service": "Booking management and reservations"
			}
		}`))
	}).Methods("GET")

	// CORS configuration - Only the API Gateway should handle CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(router)

	// Start server
	port := getEnv("PORT", "8000")
	log.Printf("🚀 API Gateway started on port %s", port)
	log.Printf("   Auth Service: %s", authServiceURL)
	log.Printf("   Building Service: %s", buildingServiceURL)
	log.Printf("   Booking Service: %s", bookingServiceURL)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// createProxyHandler creates a reverse proxy handler for the given service URL
func createProxyHandler(serviceURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the service URL
		target, err := url.Parse(serviceURL)
		if err != nil {
			log.Printf("Error parsing service URL: %v", err)
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Modify the request
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = target.Host

		// Custom error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error for %s: %v", r.URL.Path, err)
			
			// Check if we can still write to response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			
			serviceName := getServiceName(r.URL.Path)
			io.WriteString(w, `{"success": false, "error": "Service `+serviceName+` is currently unavailable"}`)
		}

		// Serve the request
		proxy.ServeHTTP(w, r)
	}
}

// getServiceName extracts service name from path
func getServiceName(path string) string {
	if strings.HasPrefix(path, "/api/auth") {
		return "auth"
	} else if strings.HasPrefix(path, "/api/buildings") {
		return "building"
	} else if strings.HasPrefix(path, "/api/bookings") {
		return "booking"
	}
	return "unknown"
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
