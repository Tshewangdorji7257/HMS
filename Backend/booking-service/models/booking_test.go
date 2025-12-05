package models

import (
	"testing"
	"time"
)

func TestBookingModel(t *testing.T) {
	now := time.Now()
	booking := Booking{
		ID:           "booking-123",
		UserID:       "user-456",
		UserName:     "Test User",
		BuildingID:   "building-789",
		BuildingName: "North Wing",
		RoomID:       "room-101",
		RoomNumber:   "101",
		BedID:        "bed-001",
		BedNumber:    1,
		BookingDate:  now,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	
	if booking.ID != "booking-123" {
		t.Errorf("Expected ID booking-123, got %s", booking.ID)
	}
	
	if booking.UserName != "Test User" {
		t.Errorf("Expected UserName Test User, got %s", booking.UserName)
	}
	
	if booking.Status != "active" {
		t.Errorf("Expected Status active, got %s", booking.Status)
	}
	
	if booking.RoomNumber != "101" {
		t.Errorf("Expected RoomNumber 101, got %s", booking.RoomNumber)
	}
}

func TestCreateBookingRequest(t *testing.T) {
	req := CreateBookingRequest{
		UserID:       "user-123",
		UserName:     "John Doe",
		BuildingID:   "building-456",
		BuildingName: "South Wing",
		RoomID:       "room-789",
		RoomNumber:   "202",
		BedID:        "bed-001",
		BedNumber:    2,
		UserEmail:    "john@example.com",
	}
	
	if req.UserID != "user-123" {
		t.Errorf("Expected UserID user-123, got %s", req.UserID)
	}
	
	if req.UserName != "John Doe" {
		t.Errorf("Expected UserName John Doe, got %s", req.UserName)
	}
	
	if req.UserEmail != "john@example.com" {
		t.Errorf("Expected UserEmail john@example.com, got %s", req.UserEmail)
	}
}

func TestBookingResponse(t *testing.T) {
	now := time.Now()
	booking := &Booking{
		ID:         "booking-xyz",
		UserName:   "Jane Smith",
		Status:     "active",
		CreatedAt:  now,
	}
	
	response := BookingResponse{
		Success: true,
		Booking: booking,
		Message: "Booking created successfully",
	}
	
	if !response.Success {
		t.Error("Expected Success to be true")
	}
	
	if response.Booking.ID != "booking-xyz" {
		t.Errorf("Expected Booking ID booking-xyz, got %s", response.Booking.ID)
	}
	
	if response.Message != "Booking created successfully" {
		t.Errorf("Expected Message 'Booking created successfully', got %s", response.Message)
	}
}

func TestBookingsResponse(t *testing.T) {
	bookings := []Booking{
		{
			ID:       "booking-1",
			UserName: "User One",
			Status:   "active",
		},
		{
			ID:       "booking-2",
			UserName: "User Two",
			Status:   "cancelled",
		},
		{
			ID:       "booking-3",
			UserName: "User Three",
			Status:   "active",
		},
	}
	
	response := BookingsResponse{
		Success:  true,
		Bookings: bookings,
	}
	
	if !response.Success {
		t.Error("Expected Success to be true")
	}
	
	if len(response.Bookings) != 3 {
		t.Errorf("Expected 3 bookings, got %d", len(response.Bookings))
	}
}

func TestBookingStatusValues(t *testing.T) {
	validStatuses := []string{"active", "cancelled", "completed"}
	
	for _, status := range validStatuses {
		booking := Booking{
			ID:     "test-booking",
			Status: status,
		}
		
		if booking.Status != status {
			t.Errorf("Expected Status %s, got %s", status, booking.Status)
		}
	}
}

func TestBookingErrorResponse(t *testing.T) {
	response := BookingResponse{
		Success: false,
		Error:   "Bed not available",
	}
	
	if response.Success {
		t.Error("Expected Success to be false")
	}
	
	if response.Error != "Bed not available" {
		t.Errorf("Expected Error 'Bed not available', got %s", response.Error)
	}
	
	if response.Booking != nil {
		t.Error("Expected Booking to be nil on error")
	}
}

func TestBookingsErrorResponse(t *testing.T) {
	response := BookingsResponse{
		Success: false,
		Error:   "Database connection failed",
	}
	
	if response.Success {
		t.Error("Expected Success to be false")
	}
	
	if response.Error != "Database connection failed" {
		t.Errorf("Expected Error 'Database connection failed', got %s", response.Error)
	}
	
	if len(response.Bookings) != 0 {
		t.Errorf("Expected 0 bookings on error, got %d", len(response.Bookings))
	}
}
