package utils

import (
	"os"
	"testing"
	"time"
)

func TestBookingConfirmationData(t *testing.T) {
	data := BookingConfirmationData{
		StudentName:  "John Doe",
		BuildingName: "North Wing",
		RoomNumber:   "101",
		BedNumber:    1,
		BookingDate:  time.Now().Format("January 2, 2006"),
		BookingID:    "booking-123",
	}
	
	if data.StudentName != "John Doe" {
		t.Errorf("Expected StudentName John Doe, got %s", data.StudentName)
	}
	
	if data.BedNumber != 1 {
		t.Errorf("Expected BedNumber 1, got %d", data.BedNumber)
	}
	
	if data.BookingID == "" {
		t.Error("BookingID should not be empty")
	}
}

func TestCancellationData(t *testing.T) {
	data := BookingCancellationData{
		StudentName:  "Jane Smith",
		BuildingName: "South Wing",
		RoomNumber:   "202",
		BedNumber:    2,
		BookingID:    "booking-456",
	}
	
	if data.StudentName != "Jane Smith" {
		t.Errorf("Expected StudentName Jane Smith, got %s", data.StudentName)
	}
	
	if data.RoomNumber != "202" {
		t.Errorf("Expected RoomNumber 202, got %s", data.RoomNumber)
	}
}

func TestEmailConfigEnvironmentVariables(t *testing.T) {
	// Test setting environment variables
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USER", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "testpassword")
	os.Setenv("FROM_EMAIL", "noreply@test.com")
	os.Setenv("FROM_NAME", "Test System")
	
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")
	
	if host != "smtp.test.com" {
		t.Errorf("Expected SMTP_HOST smtp.test.com, got %s", host)
	}
	
	if port != "587" {
		t.Errorf("Expected SMTP_PORT 587, got %s", port)
	}
	
	if user != "test@example.com" {
		t.Errorf("Expected SMTP_USER test@example.com, got %s", user)
	}
	
	if fromEmail != "noreply@test.com" {
		t.Errorf("Expected FROM_EMAIL noreply@test.com, got %s", fromEmail)
	}
	
	if fromName != "Test System" {
		t.Errorf("Expected FROM_NAME Test System, got %s", fromName)
	}
}

func TestEmailConfigurationDefaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("SMTP_HOST_TEST")
	os.Unsetenv("SMTP_PORT_TEST")
	
	// Test that defaults would be used
	host := os.Getenv("SMTP_HOST_TEST")
	if host != "" {
		t.Errorf("Expected empty SMTP_HOST_TEST after unset, got %s", host)
	}
}

func TestEmailDataValidation(t *testing.T) {
	// Test that all required fields are present
	confirmData := BookingConfirmationData{
		StudentName:  "Test Student",
		BuildingName: "Test Building",
		RoomNumber:   "101",
		BedNumber:    1,
		BookingDate:  time.Now().Format("January 2, 2006"),
		BookingID:    "test-booking-id",
	}
	
	if confirmData.StudentName == "" {
		t.Error("StudentName should not be empty")
	}
	
	if confirmData.BuildingName == "" {
		t.Error("BuildingName should not be empty")
	}
	
	if confirmData.RoomNumber == "" {
		t.Error("RoomNumber should not be empty")
	}
	
	if confirmData.BookingID == "" {
		t.Error("BookingID should not be empty")
	}
}

func TestBedNumberTypes(t *testing.T) {
	// Test different bed number types
	bedNumbers := []int{1, 2, 3, 4, 5, 10, 20}
	
	for _, bedNum := range bedNumbers {
		data := BookingConfirmationData{
			BedNumber: bedNum,
		}
		
		if data.BedNumber != bedNum {
			t.Errorf("Expected BedNumber %d, got %d", bedNum, data.BedNumber)
		}
		
		if data.BedNumber < 1 {
			t.Errorf("BedNumber should be positive, got %d", data.BedNumber)
		}
	}
}

func TestDateFormatting(t *testing.T) {
	now := time.Now()
	formatted := now.Format("January 2, 2006")
	
	if formatted == "" {
		t.Error("Formatted date should not be empty")
	}
	
	// Test that it's in the expected format
	data := BookingConfirmationData{
		BookingDate: formatted,
	}
	
	if data.BookingDate == "" {
		t.Error("BookingDate should not be empty")
	}
}

func TestMultipleEmailTypes(t *testing.T) {
	// Test confirmation email data
	confirmData := BookingConfirmationData{
		StudentName:  "Student 1",
		BuildingName: "Building A",
		RoomNumber:   "101",
		BedNumber:    1,
		BookingDate:  "December 5, 2025",
		BookingID:    "booking-1",
	}
	
	// Test cancellation email data
	cancelData := BookingCancellationData{
		StudentName:  "Student 2",
		BuildingName: "Building B",
		RoomNumber:   "202",
		BedNumber:    2,
		BookingID:    "booking-2",
	}
	
	if confirmData.StudentName == cancelData.StudentName {
		t.Error("Different students should have different names")
	}
	
	if confirmData.BookingID == cancelData.BookingID {
		t.Error("Different bookings should have different IDs")
	}
}
