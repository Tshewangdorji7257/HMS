package handlers

import (
	"building-service/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestBuildingsResponseStructure(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.BuildingsResponse{
		Success: true,
		Buildings: []models.BuildingWithRooms{
			{
				Building: models.Building{
					ID:            "building-1",
					Name:          "North Wing",
					Description:   "Modern building",
					TotalRooms:    50,
					TotalBeds:     200,
					AvailableBeds: 150,
				},
			},
		},
	}
	
	respondJSON(w, http.StatusOK, response)
	
	var decoded models.BuildingsResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if !decoded.Success {
		t.Error("Expected success to be true")
	}
	
	if len(decoded.Buildings) != 1 {
		t.Errorf("Expected 1 building, got %d", len(decoded.Buildings))
	}
}

func TestBuildingResponseError(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.BuildingResponse{
		Success: false,
		Error:   "Building not found",
	}
	
	respondJSON(w, http.StatusNotFound, response)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	
	var decoded models.BuildingResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if decoded.Success {
		t.Error("Expected success to be false")
	}
	
	if decoded.Error != "Building not found" {
		t.Errorf("Expected error message 'Building not found', got %s", decoded.Error)
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

func TestJSONEncoding(t *testing.T) {
	w := httptest.NewRecorder()
	
	building := models.Building{
		ID:          "test-building",
		Name:        "Test Building",
		Description: "Test Description",
		Amenities:   []string{"WiFi", "AC", "Gym"},
	}
	
	respondJSON(w, http.StatusOK, building)
	
	var decoded models.Building
	err := json.NewDecoder(w.Body).Decode(&decoded)
	
	if err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
	}
	
	if decoded.ID != building.ID {
		t.Errorf("Expected ID %s, got %s", building.ID, decoded.ID)
	}
	
	if len(decoded.Amenities) != len(building.Amenities) {
		t.Errorf("Expected %d amenities, got %d", len(building.Amenities), len(decoded.Amenities))
	}
}

func TestEmptyBuildingsResponse(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.BuildingsResponse{
		Success:   true,
		Buildings: []models.BuildingWithRooms{},
	}
	
	respondJSON(w, http.StatusOK, response)
	
	var decoded models.BuildingsResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if !decoded.Success {
		t.Error("Expected success to be true")
	}
	
	if len(decoded.Buildings) != 0 {
		t.Errorf("Expected 0 buildings, got %d", len(decoded.Buildings))
	}
}

func TestInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	
	response := models.BuildingsResponse{
		Success: false,
		Error:   "Internal server error",
	}
	
	respondJSON(w, http.StatusInternalServerError, response)
	
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestMultipleBuildingsResponse(t *testing.T) {
	w := httptest.NewRecorder()
	
	buildings := []models.BuildingWithRooms{
		{Building: models.Building{ID: "1", Name: "Building 1"}},
		{Building: models.Building{ID: "2", Name: "Building 2"}},
		{Building: models.Building{ID: "3", Name: "Building 3"}},
	}
	
	response := models.BuildingsResponse{
		Success:   true,
		Buildings: buildings,
	}
	
	respondJSON(w, http.StatusOK, response)
	
	var decoded models.BuildingsResponse
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if len(decoded.Buildings) != 3 {
		t.Errorf("Expected 3 buildings, got %d", len(decoded.Buildings))
	}
}

func TestBuildingWithAmenities(t *testing.T) {
	building := models.Building{
		ID:          "test-building",
		Name:        "Test Building",
		Description: "Test Description",
		Amenities:   []string{"WiFi", "AC", "Gym", "Laundry"},
	}
	
	if len(building.Amenities) != 4 {
		t.Errorf("Expected 4 amenities, got %d", len(building.Amenities))
	}
	
	w := httptest.NewRecorder()
	respondJSON(w, http.StatusOK, building)
	
	var decoded models.Building
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if len(decoded.Amenities) != 4 {
		t.Errorf("Expected 4 amenities after encoding, got %d", len(decoded.Amenities))
	}
}

func TestRoomWithBeds(t *testing.T) {
	room := models.Room{
		ID:            "room-1",
		Number:        "101",
		Type:          "double",
		TotalBeds:     2,
		AvailableBeds: 1,
	}
	
	if room.TotalBeds != 2 {
		t.Errorf("Expected 2 total beds, got %d", room.TotalBeds)
	}
	
	w := httptest.NewRecorder()
	respondJSON(w, http.StatusOK, room)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestBedOccupancy(t *testing.T) {
	userID := "user-123"
	userName := "Test User"
	
	bed := models.Bed{
		ID:             "bed-1",
		Number:         1,
		IsOccupied:     true,
		OccupiedBy:     &userID,
		OccupiedByName: &userName,
	}
	
	if !bed.IsOccupied {
		t.Error("Expected bed to be occupied")
	}
	
	if bed.OccupiedBy == nil {
		t.Error("Expected OccupiedBy to be set")
	}
	
	w := httptest.NewRecorder()
	respondJSON(w, http.StatusOK, bed)
	
	var decoded models.Bed
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if !decoded.IsOccupied {
		t.Error("Expected decoded bed to be occupied")
	}
}

func TestBuildingStatistics(t *testing.T) {
	building := models.Building{
		ID:            "building-stats",
		Name:          "Stats Building",
		TotalRooms:    100,
		TotalBeds:     400,
		AvailableBeds: 250,
	}
	
	occupancyRate := float64(building.TotalBeds-building.AvailableBeds) / float64(building.TotalBeds) * 100
	
	if occupancyRate < 0 || occupancyRate > 100 {
		t.Errorf("Invalid occupancy rate: %.2f", occupancyRate)
	}
	
	if building.AvailableBeds > building.TotalBeds {
		t.Error("Available beds cannot exceed total beds")
	}
}

func TestStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"Created", http.StatusCreated},
		{"BadRequest", http.StatusBadRequest},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
		{"Conflict", http.StatusConflict},
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

func TestComplexBuildingStructure(t *testing.T) {
	buildingWithRooms := models.BuildingWithRooms{
		Building: models.Building{
			ID:            "complex-building",
			Name:          "Complex Building",
			TotalRooms:    10,
			TotalBeds:     40,
			AvailableBeds: 30,
		},
		Rooms: []models.RoomWithBeds{
			{
				Room: models.Room{
					ID:            "room-1",
					Number:        "101",
					TotalBeds:     4,
					AvailableBeds: 3,
				},
				Beds: []models.Bed{
					{ID: "bed-1", Number: 1, IsOccupied: false},
					{ID: "bed-2", Number: 2, IsOccupied: true},
					{ID: "bed-3", Number: 3, IsOccupied: false},
					{ID: "bed-4", Number: 4, IsOccupied: false},
				},
			},
		},
	}
	
	if len(buildingWithRooms.Rooms) != 1 {
		t.Errorf("Expected 1 room, got %d", len(buildingWithRooms.Rooms))
	}
	
	if len(buildingWithRooms.Rooms[0].Beds) != 4 {
		t.Errorf("Expected 4 beds, got %d", len(buildingWithRooms.Rooms[0].Beds))
	}
	
	w := httptest.NewRecorder()
	respondJSON(w, http.StatusOK, buildingWithRooms)
	
	var decoded models.BuildingWithRooms
	json.NewDecoder(w.Body).Decode(&decoded)
	
	if decoded.Building.ID != "complex-building" {
		t.Error("Building ID mismatch after decode")
	}
}
