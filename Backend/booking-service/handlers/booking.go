package handlers

import (
	"booking-service/database"
	"booking-service/models"
	"booking-service/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateBooking creates a new booking
func CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, models.BookingResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// Validate request
	if req.UserID == "" || req.BedID == "" {
		respondJSON(w, http.StatusBadRequest, models.BookingResponse{
			Success: false,
			Error:   "User ID and Bed ID are required",
		})
		return
	}

	// Check if user already has an active booking
	var existingID string
	err := database.DB.QueryRow(
		"SELECT id FROM bookings WHERE user_id = $1 AND status = 'active'",
		req.UserID,
	).Scan(&existingID)

	if err == nil {
		respondJSON(w, http.StatusConflict, models.BookingResponse{
			Success: false,
			Error:   "You already have an active booking. Cancel it first to book a new bed.",
		})
		return
	}

	// Check if bed is already occupied
	err = database.DB.QueryRow(
		"SELECT id FROM bookings WHERE bed_id = $1 AND status = 'active'",
		req.BedID,
	).Scan(&existingID)

	if err == nil {
		respondJSON(w, http.StatusConflict, models.BookingResponse{
			Success: false,
			Error:   "This bed is already occupied",
		})
		return
	}

	// Create booking
	booking := &models.Booking{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		UserName:     req.UserName,
		BuildingID:   req.BuildingID,
		BuildingName: req.BuildingName,
		RoomID:       req.RoomID,
		RoomNumber:   req.RoomNumber,
		BedID:        req.BedID,
		BedNumber:    req.BedNumber,
		BookingDate:  time.Now(),
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = database.DB.Exec(`
		INSERT INTO bookings (
			id, user_id, user_name, building_id, building_name, 
			room_id, room_number, bed_id, bed_number, booking_date, 
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`,
		booking.ID, booking.UserID, booking.UserName, booking.BuildingID, booking.BuildingName,
		booking.RoomID, booking.RoomNumber, booking.BedID, booking.BedNumber, booking.BookingDate,
		booking.Status, booking.CreatedAt, booking.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating booking: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.BookingResponse{
			Success: false,
			Error:   "Failed to create booking",
		})
		return
	}

	// Update bed occupancy in building service
	if err := updateBedOccupancy(req.BedID, true, req.UserID, req.UserName); err != nil {
		log.Printf("Error updating bed occupancy: %v", err)
		// Rollback booking if bed update fails
		database.DB.Exec("DELETE FROM bookings WHERE id = $1", booking.ID)
		respondJSON(w, http.StatusInternalServerError, models.BookingResponse{
			Success: false,
			Error:   "Failed to update bed occupancy",
		})
		return
	}

	respondJSON(w, http.StatusCreated, models.BookingResponse{
		Success: true,
		Message: "Booking created successfully",
		Booking: booking,
	})
}

// GetAllBookings returns all bookings
func GetAllBookings(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, user_name, building_id, building_name, 
		       room_id, room_number, bed_id, bed_number, booking_date, 
		       status, created_at, updated_at 
		FROM bookings ORDER BY booking_date DESC
	`)
	if err != nil {
		log.Printf("Error fetching bookings: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.BookingsResponse{
			Success: false,
			Error:   "Failed to fetch bookings",
		})
		return
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.UserName, &booking.BuildingID, &booking.BuildingName,
			&booking.RoomID, &booking.RoomNumber, &booking.BedID, &booking.BedNumber, &booking.BookingDate,
			&booking.Status, &booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning booking: %v", err)
			continue
		}
		bookings = append(bookings, booking)
	}

	respondJSON(w, http.StatusOK, models.BookingsResponse{
		Success:  true,
		Bookings: bookings,
	})
}

// GetBookingsByUserID returns all bookings for a specific user
func GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	rows, err := database.DB.Query(`
		SELECT id, user_id, user_name, building_id, building_name, 
		       room_id, room_number, bed_id, bed_number, booking_date, 
		       status, created_at, updated_at 
		FROM bookings WHERE user_id = $1 ORDER BY booking_date DESC
	`, userID)
	if err != nil {
		log.Printf("Error fetching user bookings: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.BookingsResponse{
			Success: false,
			Error:   "Failed to fetch bookings",
		})
		return
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.UserName, &booking.BuildingID, &booking.BuildingName,
			&booking.RoomID, &booking.RoomNumber, &booking.BedID, &booking.BedNumber, &booking.BookingDate,
			&booking.Status, &booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			continue
		}
		bookings = append(bookings, booking)
	}

	respondJSON(w, http.StatusOK, models.BookingsResponse{
		Success:  true,
		Bookings: bookings,
	})
}

// GetBookingByID returns a specific booking
func GetBookingByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID := vars["id"]

	var booking models.Booking
	err := database.DB.QueryRow(`
		SELECT id, user_id, user_name, building_id, building_name, 
		       room_id, room_number, bed_id, bed_number, booking_date, 
		       status, created_at, updated_at 
		FROM bookings WHERE id = $1
	`, bookingID).Scan(
		&booking.ID, &booking.UserID, &booking.UserName, &booking.BuildingID, &booking.BuildingName,
		&booking.RoomID, &booking.RoomNumber, &booking.BedID, &booking.BedNumber, &booking.BookingDate,
		&booking.Status, &booking.CreatedAt, &booking.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, models.BookingResponse{
			Success: false,
			Error:   "Booking not found",
		})
		return
	} else if err != nil {
		log.Printf("Error fetching booking: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.BookingResponse{
			Success: false,
			Error:   "Failed to fetch booking",
		})
		return
	}

	respondJSON(w, http.StatusOK, models.BookingResponse{
		Success: true,
		Booking: &booking,
	})
}

// CancelBooking cancels a booking
func CancelBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID := vars["id"]

	// Get booking details before cancellation
	var booking models.Booking
	err := database.DB.QueryRow(`
		SELECT id, user_id, user_name, building_id, building_name, 
		       room_id, room_number, bed_id, bed_number, booking_date, 
		       status, created_at, updated_at 
		FROM bookings WHERE id = $1
	`, bookingID).Scan(
		&booking.ID, &booking.UserID, &booking.UserName, &booking.BuildingID, &booking.BuildingName,
		&booking.RoomID, &booking.RoomNumber, &booking.BedID, &booking.BedNumber, &booking.BookingDate,
		&booking.Status, &booking.CreatedAt, &booking.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, models.BookingResponse{
			Success: false,
			Error:   "Booking not found",
		})
		return
	} else if err != nil {
		log.Printf("Error fetching booking: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.BookingResponse{
			Success: false,
			Error:   "Failed to fetch booking",
		})
		return
	}

	if booking.Status == "cancelled" {
		respondJSON(w, http.StatusBadRequest, models.BookingResponse{
			Success: false,
			Error:   "Booking is already cancelled",
		})
		return
	}

	// Update booking status to cancelled
	_, err = database.DB.Exec(
		"UPDATE bookings SET status = 'cancelled', updated_at = $1 WHERE id = $2",
		time.Now(), bookingID,
	)
	if err != nil {
		log.Printf("Error cancelling booking: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.BookingResponse{
			Success: false,
			Error:   "Failed to cancel booking",
		})
		return
	}

	// Update bed occupancy in building service
	if err := updateBedOccupancy(booking.BedID, false, "", ""); err != nil {
		log.Printf("Error updating bed occupancy: %v", err)
	}

	booking.Status = "cancelled"
	booking.UpdatedAt = time.Now()

	respondJSON(w, http.StatusOK, models.BookingResponse{
		Success: true,
		Message: "Booking cancelled successfully",
		Booking: &booking,
	})
}

// Helper function to update bed occupancy in building service
func updateBedOccupancy(bedID string, isOccupied bool, occupiedBy, occupiedByName string) error {
	buildingServiceURL := utils.GetBuildingServiceURL()

	var payload map[string]interface{}
	if isOccupied {
		payload = map[string]interface{}{
			"is_occupied":      true,
			"occupied_by":      occupiedBy,
			"occupied_by_name": occupiedByName,
		}
	} else {
		payload = map[string]interface{}{
			"is_occupied":      false,
			"occupied_by":      nil,
			"occupied_by_name": nil,
		}
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/buildings/beds/%s/occupancy", buildingServiceURL, bedID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update bed occupancy, status code: %d", resp.StatusCode)
	}

	return nil
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
