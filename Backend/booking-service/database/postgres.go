package database

import (
	"database/sql"
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

	return nil
}

// createTables creates necessary database tables
func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS bookings (
		id VARCHAR(255) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		user_name VARCHAR(255) NOT NULL,
		building_id VARCHAR(255) NOT NULL,
		building_name VARCHAR(255) NOT NULL,
		room_id VARCHAR(255) NOT NULL,
		room_number VARCHAR(50) NOT NULL,
		bed_id VARCHAR(255) NOT NULL,
		bed_number INT NOT NULL,
		booking_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(50) DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_bookings_user ON bookings(user_id);
	CREATE INDEX IF NOT EXISTS idx_bookings_bed ON bookings(bed_id);
	CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
	CREATE INDEX IF NOT EXISTS idx_bookings_building ON bookings(building_id);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("✅ Database tables created/verified successfully")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
