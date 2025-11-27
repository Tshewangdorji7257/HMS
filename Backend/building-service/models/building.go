package models

import "time"

// Building represents a hostel building
type Building struct {
	ID            string    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	TotalRooms    int       `json:"total_rooms" db:"total_rooms"`
	TotalBeds     int       `json:"total_beds" db:"total_beds"`
	AvailableBeds int       `json:"available_beds" db:"available_beds"`
	Amenities     []string  `json:"amenities" db:"amenities"`
	Image         string    `json:"image" db:"image"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Room represents a room in a building
type Room struct {
	ID            string    `json:"id" db:"id"`
	BuildingID    string    `json:"building_id" db:"building_id"`
	Number        string    `json:"number" db:"number"`
	Type          string    `json:"type" db:"type"` // "single", "double", "triple", "quad"
	TotalBeds     int       `json:"total_beds" db:"total_beds"`
	AvailableBeds int       `json:"available_beds" db:"available_beds"`
	Amenities     []string  `json:"amenities" db:"amenities"`
	Price         float64   `json:"price" db:"price"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Bed represents a bed in a room
type Bed struct {
	ID             string  `json:"id" db:"id"`
	RoomID         string  `json:"room_id" db:"room_id"`
	Number         int     `json:"number" db:"number"`
	IsOccupied     bool    `json:"is_occupied" db:"is_occupied"`
	OccupiedBy     *string `json:"occupied_by,omitempty" db:"occupied_by"`
	OccupiedByName *string `json:"occupied_by_name,omitempty" db:"occupied_by_name"`
}

// BuildingWithRooms represents a building with its rooms
type BuildingWithRooms struct {
	Building
	Rooms []RoomWithBeds `json:"rooms"`
}

// RoomWithBeds represents a room with its beds
type RoomWithBeds struct {
	Room
	Beds []Bed `json:"beds"`
}

// BuildingResponse represents API response for buildings
type BuildingResponse struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message,omitempty"`
	Building interface{} `json:"building,omitempty"`
	Error    string      `json:"error,omitempty"`
}

// BuildingsResponse represents API response for multiple buildings
type BuildingsResponse struct {
	Success   bool                `json:"success"`
	Buildings []BuildingWithRooms `json:"buildings,omitempty"`
	Error     string              `json:"error,omitempty"`
}
