package models

import "time"

// Booking represents a room booking
type Booking struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	UserName     string    `json:"user_name" db:"user_name"`
	BuildingID   string    `json:"building_id" db:"building_id"`
	BuildingName string    `json:"building_name" db:"building_name"`
	RoomID       string    `json:"room_id" db:"room_id"`
	RoomNumber   string    `json:"room_number" db:"room_number"`
	BedID        string    `json:"bed_id" db:"bed_id"`
	BedNumber    int       `json:"bed_number" db:"bed_number"`
	BookingDate  time.Time `json:"booking_date" db:"booking_date"`
	Status       string    `json:"status" db:"status"` // "active" or "cancelled"
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateBookingRequest represents a booking creation request
type CreateBookingRequest struct {
	UserID       string `json:"user_id" binding:"required"`
	UserName     string `json:"user_name" binding:"required"`
	BuildingID   string `json:"building_id" binding:"required"`
	BuildingName string `json:"building_name" binding:"required"`
	RoomID       string `json:"room_id" binding:"required"`
	RoomNumber   string `json:"room_number" binding:"required"`
	BedID        string `json:"bed_id" binding:"required"`
	BedNumber    int    `json:"bed_number" binding:"required"`
}

// BookingResponse represents API response for booking
type BookingResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Booking *Booking `json:"booking,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// BookingsResponse represents API response for multiple bookings
type BookingsResponse struct {
	Success  bool      `json:"success"`
	Bookings []Booking `json:"bookings,omitempty"`
	Error    string    `json:"error,omitempty"`
}
