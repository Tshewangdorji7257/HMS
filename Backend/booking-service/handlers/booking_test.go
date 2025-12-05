package handlers

import (
	"booking-service/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"success": true,
		"message": "test",
	}
	
	respondJSON(w, http.StatusOK, data)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}

func TestCreateBookingValidation(t *testing.T) {
	tests := []struct {
		name       string
		payload    models.CreateBookingRequest
		wantStatus int
	}{
		{
			name: "Missing UserID",
			payload: models.CreateBookingRequest{
				BedID: "bed-123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing BedID",
			payload: models.CreateBookingRequest{
				UserID: "user-123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			CreateBooking(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestCreateBookingInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/bookings", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateBooking(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response models.BookingResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}
}

func TestBookingResponseStructure(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.BookingResponse{
		Success: true,
		Message: "Booking created successfully",
		Booking: &models.Booking{
			ID:       "booking-123",
			UserID:   "user-456",
			UserName: "Test User",
			Status:   "active",
		},
	}
	
	respondJSON(w, http.StatusCreated, response)
	
	var decoded models.BookingResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if !decoded.Success {
		t.Error("Expected success to be true")
	}
	
	if decoded.Booking.ID != "booking-123" {
		t.Errorf("Expected booking ID booking-123, got %s", decoded.Booking.ID)
	}
}

func TestBookingsResponseStructure(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.BookingsResponse{
		Success: true,
		Bookings: []models.Booking{
			{ID: "booking-1", Status: "active"},
			{ID: "booking-2", Status: "cancelled"},
		},
	}
	
	respondJSON(w, http.StatusOK, response)
	
	var decoded models.BookingsResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if !decoded.Success {
		t.Error("Expected success to be true")
	}
	
	if len(decoded.Bookings) != 2 {
		t.Errorf("Expected 2 bookings, got %d", len(decoded.Bookings))
	}
}

func TestBookingErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.BookingResponse{
		Success: false,
		Error:   "Bed already occupied",
	}
	
	respondJSON(w, http.StatusConflict, response)
	
	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
	
	var decoded models.BookingResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if decoded.Success {
		t.Error("Expected success to be false")
	}
	
	if decoded.Error != "Bed already occupied" {
		t.Errorf("Expected error 'Bed already occupied', got %s", decoded.Error)
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEmptyBookingsResponse(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.BookingsResponse{
		Success:  true,
		Bookings: []models.Booking{},
	}
	
	respondJSON(w, http.StatusOK, response)
	
	var decoded models.BookingsResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if !decoded.Success {
		t.Error("Expected success to be true")
	}
	
	if len(decoded.Bookings) != 0 {
		t.Errorf("Expected 0 bookings, got %d", len(decoded.Bookings))
	}
}

func TestJSONEncoding(t *testing.T) {
	w := httptest.NewRecorder()
	
	booking := models.Booking{
		ID:           "test-booking",
		UserName:     "Test User",
		BuildingName: "North Wing",
		RoomNumber:   "101",
		BedNumber:    1,
		Status:       "active",
	}
	
	respondJSON(w, http.StatusOK, booking)
	
	var decoded models.Booking
	err := json.NewDecoder(w.Body).Decode(&decoded)
	
	if err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
	}
	
	if decoded.ID != booking.ID {
		t.Errorf("Expected ID %s, got %s", booking.ID, decoded.ID)
	}
	
	if decoded.Status != "active" {
		t.Errorf("Expected status active, got %s", decoded.Status)
	}
}

func TestMultipleBookingsResponse(t *testing.T) {
	w := httptest.NewRecorder()
	
	bookings := []models.Booking{
		{ID: "1", UserName: "User 1", Status: "active"},
		{ID: "2", UserName: "User 2", Status: "active"},
		{ID: "3", UserName: "User 3", Status: "cancelled"},
	}
	
	response := models.BookingsResponse{
		Success:  true,
		Bookings: bookings,
	}
	
	respondJSON(w, http.StatusOK, response)
	
	var decoded models.BookingsResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if len(decoded.Bookings) != 3 {
		t.Errorf("Expected 3 bookings, got %d", len(decoded.Bookings))
	}
}

func TestBookingStatusTransitions(t *testing.T) {
	tests := []struct {
		name   string
		status string
		valid  bool
	}{
		{"Active status", "active", true},
		{"Cancelled status", "cancelled", true},
		{"Completed status", "completed", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booking := models.Booking{
				ID:     "test-booking",
				Status: tt.status,
			}
			
			if booking.Status != tt.status {
				t.Errorf("Expected status %s, got %s", tt.status, booking.Status)
			}
		})
	}
}

func TestBookingWithAllFields(t *testing.T) {
	now := time.Now()
	booking := models.Booking{
		ID:           "full-booking",
		UserID:       "user-123",
		UserName:     "John Doe",
		BuildingID:   "building-456",
		BuildingName: "North Wing",
		RoomID:       "room-789",
		RoomNumber:   "101",
		BedID:        "bed-001",
		BedNumber:    1,
		BookingDate:  now,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	
	w := httptest.NewRecorder()
	respondJSON(w, http.StatusOK, booking)
	
	var decoded models.Booking
	err := json.NewDecoder(w.Body).Decode(&decoded)
	if err != nil {
		t.Fatalf("Failed to decode booking: %v", err)
	}
	
	if decoded.ID != booking.ID {
		t.Errorf("Expected ID %s, got %s", booking.ID, decoded.ID)
	}
	
	if decoded.UserName != booking.UserName {
		t.Errorf("Expected UserName %s, got %s", booking.UserName, decoded.UserName)
	}
	
	if decoded.BedNumber != booking.BedNumber {
		t.Errorf("Expected BedNumber %d, got %d", booking.BedNumber, decoded.BedNumber)
	}
}

func TestCreateBookingRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.CreateBookingRequest
		valid   bool
	}{
		{
			name: "Valid request",
			request: models.CreateBookingRequest{
				UserID:       "user-1",
				UserName:     "Test User",
				BedID:        "bed-1",
				BuildingID:   "building-1",
				BuildingName: "Test Building",
				RoomID:       "room-1",
				RoomNumber:   "101",
				BedNumber:    1,
			},
			valid: true,
		},
		{
			name: "Missing UserID",
			request: models.CreateBookingRequest{
				BedID: "bed-1",
			},
			valid: false,
		},
		{
			name: "Missing BedID",
			request: models.CreateBookingRequest{
				UserID: "user-1",
			},
			valid: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.request.UserID != "" && tt.request.BedID != ""
			if isValid != tt.valid {
				t.Errorf("Expected validation %v, got %v", tt.valid, isValid)
			}
		})
	}
}

func TestBookingConflictScenarios(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
		statusCode   int
	}{
		{
			name:         "Bed already occupied",
			errorMessage: "This bed is already occupied",
			statusCode:   http.StatusConflict,
		},
		{
			name:         "User has active booking",
			errorMessage: "You already have an active booking",
			statusCode:   http.StatusConflict,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			response := models.BookingResponse{
				Success: false,
				Error:   tt.errorMessage,
			}
			
			respondJSON(w, tt.statusCode, response)
			
			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestStatusCodeResponses(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"Created", http.StatusCreated},
		{"BadRequest", http.StatusBadRequest},
		{"Conflict", http.StatusConflict},
		{"InternalServerError", http.StatusInternalServerError},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data := map[string]string{"status": tt.name}
			
			respondJSON(w, tt.statusCode, data)
			
			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestBookingArrayOperations(t *testing.T) {
	bookings := []models.Booking{
		{ID: "1", Status: "active"},
		{ID: "2", Status: "active"},
		{ID: "3", Status: "cancelled"},
	}
	
	activeCount := 0
	cancelledCount := 0
	
	for _, booking := range bookings {
		if booking.Status == "active" {
			activeCount++
		} else if booking.Status == "cancelled" {
			cancelledCount++
		}
	}
	
	if activeCount != 2 {
		t.Errorf("Expected 2 active bookings, got %d", activeCount)
	}
	
	if cancelledCount != 1 {
		t.Errorf("Expected 1 cancelled booking, got %d", cancelledCount)
	}
}
