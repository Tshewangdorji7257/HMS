package handlers

import (
	"building-service/database"
	"building-service/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// GetAllBuildings returns all buildings with their rooms and beds
func GetAllBuildings(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, name, description, total_rooms, total_beds, available_beds, 
		       COALESCE(amenities, '[]'::jsonb), COALESCE(image, ''), created_at, updated_at 
		FROM buildings ORDER BY name
	`)
	if err != nil {
		log.Printf("Error fetching buildings: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.BuildingsResponse{
			Success: false,
			Error:   "Failed to fetch buildings",
		})
		return
	}
	defer rows.Close()

	var buildings []models.BuildingWithRooms

	for rows.Next() {
		var building models.Building
		var amenitiesJSON []byte

		err := rows.Scan(
			&building.ID, &building.Name, &building.Description,
			&building.TotalRooms, &building.TotalBeds, &building.AvailableBeds,
			&amenitiesJSON, &building.Image, &building.CreatedAt, &building.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning building: %v", err)
			continue
		}

		// Parse amenities JSON
		if err := json.Unmarshal(amenitiesJSON, &building.Amenities); err != nil {
			building.Amenities = []string{}
		}

		// Get rooms for this building
		rooms, err := getRoomsForBuilding(building.ID)
		if err != nil {
			log.Printf("Error fetching rooms for building %s: %v", building.ID, err)
			rooms = []models.RoomWithBeds{}
		}

		buildings = append(buildings, models.BuildingWithRooms{
			Building: building,
			Rooms:    rooms,
		})
	}

	respondJSON(w, http.StatusOK, models.BuildingsResponse{
		Success:   true,
		Buildings: buildings,
	})
}

// GetBuildingByID returns a specific building with its rooms and beds
func GetBuildingByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buildingID := vars["id"]

	var building models.Building
	var amenitiesJSON []byte

	err := database.DB.QueryRow(`
		SELECT id, name, description, total_rooms, total_beds, available_beds, 
		       COALESCE(amenities, '[]'::jsonb), COALESCE(image, ''), created_at, updated_at 
		FROM buildings WHERE id = $1
	`, buildingID).Scan(
		&building.ID, &building.Name, &building.Description,
		&building.TotalRooms, &building.TotalBeds, &building.AvailableBeds,
		&amenitiesJSON, &building.Image, &building.CreatedAt, &building.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, models.BuildingResponse{
			Success: false,
			Error:   "Building not found",
		})
		return
	} else if err != nil {
		log.Printf("Error fetching building: %v", err)
		respondJSON(w, http.StatusInternalServerError, models.BuildingResponse{
			Success: false,
			Error:   "Failed to fetch building",
		})
		return
	}

	// Parse amenities JSON
	if err := json.Unmarshal(amenitiesJSON, &building.Amenities); err != nil {
		building.Amenities = []string{}
	}

	// Get rooms for this building
	rooms, err := getRoomsForBuilding(building.ID)
	if err != nil {
		log.Printf("Error fetching rooms: %v", err)
		rooms = []models.RoomWithBeds{}
	}

	respondJSON(w, http.StatusOK, models.BuildingResponse{
		Success: true,
		Building: models.BuildingWithRooms{
			Building: building,
			Rooms:    rooms,
		},
	})
}

// GetRoomByID returns a specific room with its beds
func GetRoomByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomId"]

	var room models.Room
	var amenitiesJSON []byte

	err := database.DB.QueryRow(`
		SELECT id, building_id, number, type, total_beds, available_beds, 
		       COALESCE(amenities, '[]'::jsonb), price, created_at, updated_at 
		FROM rooms WHERE id = $1
	`, roomID).Scan(
		&room.ID, &room.BuildingID, &room.Number, &room.Type,
		&room.TotalBeds, &room.AvailableBeds,
		&amenitiesJSON, &room.Price, &room.CreatedAt, &room.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Room not found",
		})
		return
	} else if err != nil {
		log.Printf("Error fetching room: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch room",
		})
		return
	}

	// Parse amenities JSON
	if err := json.Unmarshal(amenitiesJSON, &room.Amenities); err != nil {
		room.Amenities = []string{}
	}

	// Get beds for this room
	beds, err := getBedsForRoom(room.ID)
	if err != nil {
		log.Printf("Error fetching beds: %v", err)
		beds = []models.Bed{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"room": models.RoomWithBeds{
			Room: room,
			Beds: beds,
		},
	})
}

// UpdateBedOccupancy updates bed occupancy status
func UpdateBedOccupancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bedID := vars["bedId"]

	var req struct {
		IsOccupied     bool    `json:"is_occupied"`
		OccupiedBy     *string `json:"occupied_by"`
		OccupiedByName *string `json:"occupied_by_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Update bed occupancy
	_, err := database.DB.Exec(`
		UPDATE beds 
		SET is_occupied = $1, occupied_by = $2, occupied_by_name = $3 
		WHERE id = $4
	`, req.IsOccupied, req.OccupiedBy, req.OccupiedByName, bedID)

	if err != nil {
		log.Printf("Error updating bed occupancy: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to update bed occupancy",
		})
		return
	}

	// Get room_id to update room's available_beds count
	var roomID string
	err = database.DB.QueryRow("SELECT room_id FROM beds WHERE id = $1", bedID).Scan(&roomID)
	if err == nil {
		// Update room's available beds count
		_, _ = database.DB.Exec(`
			UPDATE rooms 
			SET available_beds = (SELECT COUNT(*) FROM beds WHERE room_id = $1 AND is_occupied = false) 
			WHERE id = $1
		`, roomID)

		// Update building's available beds count
		var buildingID string
		err = database.DB.QueryRow("SELECT building_id FROM rooms WHERE id = $1", roomID).Scan(&buildingID)
		if err == nil {
			_, _ = database.DB.Exec(`
				UPDATE buildings 
				SET available_beds = (
					SELECT COUNT(*) FROM beds 
					WHERE room_id IN (SELECT id FROM rooms WHERE building_id = $1) 
					AND is_occupied = false
				) 
				WHERE id = $1
			`, buildingID)
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Bed occupancy updated successfully",
	})
}

// Helper functions
func getRoomsForBuilding(buildingID string) ([]models.RoomWithBeds, error) {
	rows, err := database.DB.Query(`
		SELECT id, building_id, number, type, total_beds, available_beds, 
		       COALESCE(amenities, '[]'::jsonb), price, created_at, updated_at 
		FROM rooms WHERE building_id = $1 ORDER BY number
	`, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.RoomWithBeds

	for rows.Next() {
		var room models.Room
		var amenitiesJSON []byte

		err := rows.Scan(
			&room.ID, &room.BuildingID, &room.Number, &room.Type,
			&room.TotalBeds, &room.AvailableBeds,
			&amenitiesJSON, &room.Price, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// Parse amenities JSON
		if err := json.Unmarshal(amenitiesJSON, &room.Amenities); err != nil {
			room.Amenities = []string{}
		}

		// Get beds for this room
		beds, err := getBedsForRoom(room.ID)
		if err != nil {
			beds = []models.Bed{}
		}

		rooms = append(rooms, models.RoomWithBeds{
			Room: room,
			Beds: beds,
		})
	}

	return rooms, nil
}

func getBedsForRoom(roomID string) ([]models.Bed, error) {
	rows, err := database.DB.Query(`
		SELECT id, room_id, number, is_occupied, occupied_by, occupied_by_name 
		FROM beds WHERE room_id = $1 ORDER BY number
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beds []models.Bed

	for rows.Next() {
		var bed models.Bed
		var occupiedBy, occupiedByName sql.NullString

		err := rows.Scan(
			&bed.ID, &bed.RoomID, &bed.Number,
			&bed.IsOccupied, &occupiedBy, &occupiedByName,
		)
		if err != nil {
			continue
		}

		if occupiedBy.Valid {
			bed.OccupiedBy = &occupiedBy.String
		}
		if occupiedByName.Valid {
			bed.OccupiedByName = &occupiedByName.String
		}

		beds = append(beds, bed)
	}

	return beds, nil
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// GetBedsByUserID returns all beds occupied by a specific user
func GetBedsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	rows, err := database.DB.Query(`
		SELECT b.id, b.room_id, b.number, b.is_occupied, b.occupied_by, b.occupied_by_name
		FROM beds b
		WHERE b.occupied_by = $1
	`, userID)
	if err != nil {
		log.Printf("Error fetching beds for user: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch beds",
		})
		return
	}
	defer rows.Close()

	var beds []models.Bed

	for rows.Next() {
		var bed models.Bed
		var occupiedBy, occupiedByName sql.NullString

		err := rows.Scan(
			&bed.ID, &bed.RoomID, &bed.Number,
			&bed.IsOccupied, &occupiedBy, &occupiedByName,
		)
		if err != nil {
			continue
		}

		if occupiedBy.Valid {
			bed.OccupiedBy = &occupiedBy.String
		}
		if occupiedByName.Valid {
			bed.OccupiedByName = &occupiedByName.String
		}

		beds = append(beds, bed)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"beds":    beds,
	})
}

// SearchBuildings searches buildings by name or amenities
func SearchBuildings(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	
	rows, err := database.DB.Query(`
		SELECT id, name, description, total_rooms, total_beds, available_beds, 
		       COALESCE(amenities, '[]'::jsonb), COALESCE(image, ''), created_at, updated_at 
		FROM buildings 
		WHERE LOWER(name) LIKE LOWER($1) OR LOWER(description) LIKE LOWER($1)
		ORDER BY name
	`, "%"+query+"%")
	
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Printf("PostgreSQL Error: %v", pqErr)
		}
		respondJSON(w, http.StatusInternalServerError, models.BuildingsResponse{
			Success: false,
			Error:   "Failed to search buildings",
		})
		return
	}
	defer rows.Close()

	var buildings []models.BuildingWithRooms

	for rows.Next() {
		var building models.Building
		var amenitiesJSON []byte

		err := rows.Scan(
			&building.ID, &building.Name, &building.Description,
			&building.TotalRooms, &building.TotalBeds, &building.AvailableBeds,
			&amenitiesJSON, &building.Image, &building.CreatedAt, &building.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning building: %v", err)
			continue
		}

		if err := json.Unmarshal(amenitiesJSON, &building.Amenities); err != nil {
			building.Amenities = []string{}
		}

		rooms, _ := getRoomsForBuilding(building.ID)
		buildings = append(buildings, models.BuildingWithRooms{
			Building: building,
			Rooms:    rooms,
		})
	}

	respondJSON(w, http.StatusOK, models.BuildingsResponse{
		Success:   true,
		Buildings: buildings,
	})
}
