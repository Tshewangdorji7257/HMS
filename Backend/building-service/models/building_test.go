package models

import (
	"testing"
	"time"
)

func TestBuildingModel(t *testing.T) {
	now := time.Now()
	building := Building{
		ID:            "building-123",
		Name:          "North Wing",
		Description:   "Modern building",
		TotalRooms:    50,
		TotalBeds:     200,
		AvailableBeds: 150,
		Amenities:     []string{"WiFi", "AC", "Laundry"},
		Image:         "https://example.com/image.jpg",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	if building.ID != "building-123" {
		t.Errorf("Expected ID building-123, got %s", building.ID)
	}
	
	if building.Name != "North Wing" {
		t.Errorf("Expected Name North Wing, got %s", building.Name)
	}
	
	if building.TotalRooms != 50 {
		t.Errorf("Expected TotalRooms 50, got %d", building.TotalRooms)
	}
	
	if len(building.Amenities) != 3 {
		t.Errorf("Expected 3 amenities, got %d", len(building.Amenities))
	}
}

func TestRoomModel(t *testing.T) {
	now := time.Now()
	room := Room{
		ID:            "room-456",
		BuildingID:    "building-123",
		Number:        "101",
		Type:          "double",
		TotalBeds:     4,
		AvailableBeds: 2,
		Amenities:     []string{"WiFi", "AC"},
		Price:         500.00,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	if room.ID != "room-456" {
		t.Errorf("Expected ID room-456, got %s", room.ID)
	}
	
	if room.Number != "101" {
		t.Errorf("Expected Number 101, got %s", room.Number)
	}
	
	if room.TotalBeds != 4 {
		t.Errorf("Expected TotalBeds 4, got %d", room.TotalBeds)
	}
	
	if room.Type != "double" {
		t.Errorf("Expected Type double, got %s", room.Type)
	}
}

func TestBedModel(t *testing.T) {
	userName := "John Doe"
	userID := "user-123"
	
	bed := Bed{
		ID:             "bed-789",
		RoomID:         "room-456",
		Number:         1,
		IsOccupied:     true,
		OccupiedBy:     &userID,
		OccupiedByName: &userName,
	}
	
	if bed.ID != "bed-789" {
		t.Errorf("Expected ID bed-789, got %s", bed.ID)
	}
	
	if bed.Number != 1 {
		t.Errorf("Expected Number 1, got %d", bed.Number)
	}
	
	if !bed.IsOccupied {
		t.Error("Expected IsOccupied to be true")
	}
	
	if bed.OccupiedBy == nil || *bed.OccupiedBy != userID {
		t.Errorf("Expected OccupiedBy %s, got %v", userID, bed.OccupiedBy)
	}
}

func TestBuildingWithRooms(t *testing.T) {
	building := Building{
		ID:   "building-1",
		Name: "Test Building",
	}
	
	rooms := []RoomWithBeds{
		{
			Room: Room{
				ID:     "room-1",
				Number: "101",
			},
			Beds: []Bed{
				{ID: "bed-1", Number: 1},
				{ID: "bed-2", Number: 2},
			},
		},
	}
	
	buildingWithRooms := BuildingWithRooms{
		Building: building,
		Rooms:    rooms,
	}
	
	if buildingWithRooms.Building.ID != "building-1" {
		t.Error("Building ID mismatch")
	}
	
	if len(buildingWithRooms.Rooms) != 1 {
		t.Errorf("Expected 1 room, got %d", len(buildingWithRooms.Rooms))
	}
	
	if len(buildingWithRooms.Rooms[0].Beds) != 2 {
		t.Errorf("Expected 2 beds, got %d", len(buildingWithRooms.Rooms[0].Beds))
	}
}

func TestBuildingsResponse(t *testing.T) {
	response := BuildingsResponse{
		Success: true,
		Buildings: []BuildingWithRooms{
			{
				Building: Building{ID: "b1", Name: "Building 1"},
			},
			{
				Building: Building{ID: "b2", Name: "Building 2"},
			},
		},
	}
	
	if !response.Success {
		t.Error("Expected Success to be true")
	}
	
	if len(response.Buildings) != 2 {
		t.Errorf("Expected 2 buildings, got %d", len(response.Buildings))
	}
}

func TestBuildingResponse(t *testing.T) {
	buildingData := BuildingWithRooms{
		Building: Building{
			ID:   "building-xyz",
			Name: "XYZ Building",
		},
	}
	
	response := BuildingResponse{
		Success:  true,
		Building: buildingData,
	}
	
	if !response.Success {
		t.Error("Expected Success to be true")
	}
	
	if response.Message != "" {
		// Message is optional
	}
}

func TestAvailableBedCount(t *testing.T) {
	room := Room{
		ID:            "room-test",
		Number:        "202",
		TotalBeds:     4,
		AvailableBeds: 2,
	}
	
	if room.AvailableBeds != 2 {
		t.Errorf("Expected AvailableBeds 2, got %d", room.AvailableBeds)
	}
	
	if room.TotalBeds != 4 {
		t.Errorf("Expected TotalBeds 4, got %d", room.TotalBeds)
	}
}

func TestBedOccupancy(t *testing.T) {
	// Test occupied bed
	userID := "user-456"
	userName := "Jane Smith"
	occupiedBed := Bed{
		ID:             "bed-1",
		RoomID:         "room-1",
		Number:         1,
		IsOccupied:     true,
		OccupiedBy:     &userID,
		OccupiedByName: &userName,
	}
	
	if !occupiedBed.IsOccupied {
		t.Error("Expected bed to be occupied")
	}
	
	if occupiedBed.OccupiedBy == nil {
		t.Error("Expected OccupiedBy to not be nil")
	}
	
	// Test available bed
	availableBed := Bed{
		ID:         "bed-2",
		RoomID:     "room-1",
		Number:     2,
		IsOccupied: false,
	}
	
	if availableBed.IsOccupied {
		t.Error("Expected bed to be available")
	}
	
	if availableBed.OccupiedBy != nil {
		t.Error("Expected OccupiedBy to be nil for available bed")
	}
}
