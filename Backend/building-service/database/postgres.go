package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the PostgreSQL database connection
func InitDB() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")

	// Create tables
	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Seed initial data
	if err = seedData(); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	return nil
}

// createTables creates necessary database tables
func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS buildings (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		total_rooms INT DEFAULT 0,
		total_beds INT DEFAULT 0,
		available_beds INT DEFAULT 0,
		amenities JSONB,
		image TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS rooms (
		id VARCHAR(255) PRIMARY KEY,
		building_id VARCHAR(255) REFERENCES buildings(id) ON DELETE CASCADE,
		number VARCHAR(50) NOT NULL,
		type VARCHAR(50) NOT NULL,
		total_beds INT NOT NULL,
		available_beds INT NOT NULL,
		amenities JSONB,
		price DECIMAL(10,2) DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(building_id, number)
	);

	CREATE TABLE IF NOT EXISTS beds (
		id VARCHAR(255) PRIMARY KEY,
		room_id VARCHAR(255) REFERENCES rooms(id) ON DELETE CASCADE,
		number INT NOT NULL,
		is_occupied BOOLEAN DEFAULT FALSE,
		occupied_by VARCHAR(255),
		occupied_by_name VARCHAR(255),
		UNIQUE(room_id, number)
	);

	CREATE INDEX IF NOT EXISTS idx_rooms_building ON rooms(building_id);
	CREATE INDEX IF NOT EXISTS idx_beds_room ON beds(room_id);
	CREATE INDEX IF NOT EXISTS idx_beds_occupied ON beds(is_occupied);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("✅ Database tables created/verified successfully")
	return nil
}

// seedData seeds initial building data
func seedData() error {
	// Check if data already exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM buildings").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("✅ Database already seeded")
		return nil
	}

	buildings := []struct {
		ID          string
		Name        string
		Description string
		Amenities   []string
		RoomCount   int
	}{
		{"bldg-1", "RK A", "Modern residence with state-of-the-art facilities", []string{"Wi-Fi", "Laundry Room", "Common Kitchen", "Study Lounge", "Recreation Room"}, 25},
		{"bldg-2", "RK B", "Cozy accommodation with a homely atmosphere", []string{"Wi-Fi", "Laundry Room", "Study Lounge", "Gym"}, 20},
		{"bldg-3", "H A", "Spacious rooms with excellent natural lighting", []string{"Wi-Fi", "Laundry Room", "Common Kitchen", "Cafeteria"}, 30},
		{"bldg-4", "H B", "Contemporary design with eco-friendly features", []string{"Wi-Fi", "Study Lounge", "Recreation Room", "Parking"}, 22},
		{"bldg-5", "H C", "Traditional architecture with modern amenities", []string{"Wi-Fi", "Laundry Room", "Gym", "Library"}, 18},
		{"bldg-6", "H D", "Quiet location perfect for studying", []string{"Wi-Fi", "Common Kitchen", "Study Lounge"}, 15},
		{"bldg-7", "H E", "Central location with easy campus access", []string{"Wi-Fi", "Laundry Room", "Recreation Room", "Parking"}, 28},
		{"bldg-8", "H F", "Newly renovated with premium facilities", []string{"Wi-Fi", "Gym", "Cafeteria", "Library"}, 24},
		{"bldg-9", "NK", "Garden view rooms with peaceful surroundings", []string{"Wi-Fi", "Laundry Room", "Common Kitchen", "Study Lounge"}, 20},
		{"bldg-10", "Lhawang", "High-rise building with panoramic views", []string{"Wi-Fi", "Recreation Room", "Gym", "Parking", "Cafeteria"}, 35},
	}

	roomTypes := []string{"single", "double", "triple", "quad"}
	bedCounts := map[string]int{"single": 1, "double": 2, "triple": 3, "quad": 4}

	for _, bldg := range buildings {
		amenitiesJSON, _ := json.Marshal(bldg.Amenities)

		totalBeds := 0
		for i := 0; i < bldg.RoomCount; i++ {
			roomType := roomTypes[i%len(roomTypes)]
			totalBeds += bedCounts[roomType]
		}

		_, err := DB.Exec(
			`INSERT INTO buildings (id, name, description, total_rooms, total_beds, available_beds, amenities, image) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			bldg.ID, bldg.Name, bldg.Description, bldg.RoomCount, totalBeds, totalBeds, amenitiesJSON, "/placeholder.svg",
		)
		if err != nil {
			return fmt.Errorf("failed to insert building %s: %w", bldg.Name, err)
		}

		// Create rooms for each building
		for i := 1; i <= bldg.RoomCount; i++ {
			roomNumber := fmt.Sprintf("%03d", i)
			roomType := roomTypes[i%len(roomTypes)]
			totalBedsInRoom := bedCounts[roomType]
			roomID := fmt.Sprintf("%s-room-%s", bldg.ID, roomNumber)

			roomAmenities := []string{"Wi-Fi", "Study Desk", "Wardrobe"}
			if roomType == "single" {
				roomAmenities = append(roomAmenities, "Private Bathroom")
			} else {
				roomAmenities = append(roomAmenities, "Shared Bathroom")
			}
			roomAmenitiesJSON, _ := json.Marshal(roomAmenities)

			_, err := DB.Exec(
				`INSERT INTO rooms (id, building_id, number, type, total_beds, available_beds, amenities, price) 
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				roomID, bldg.ID, roomNumber, roomType, totalBedsInRoom, totalBedsInRoom, roomAmenitiesJSON, 5000.00,
			)
			if err != nil {
				return fmt.Errorf("failed to insert room %s: %w", roomNumber, err)
			}

			// Create beds for each room
			for j := 1; j <= totalBedsInRoom; j++ {
				bedID := fmt.Sprintf("%s-bed-%d", roomID, j)
				_, err := DB.Exec(
					`INSERT INTO beds (id, room_id, number, is_occupied) 
					 VALUES ($1, $2, $3, $4)`,
					bedID, roomID, j, false,
				)
				if err != nil {
					return fmt.Errorf("failed to insert bed %d: %w", j, err)
				}
			}
		}
	}

	log.Println("✅ Database seeded with sample data")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
