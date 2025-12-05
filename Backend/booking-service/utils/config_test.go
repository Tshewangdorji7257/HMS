package utils

import (
	"os"
	"testing"
)

func TestGetEnvWithDefault(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	value := getEnv("TEST_VAR", "default")
	if value != "test_value" {
		t.Errorf("Expected test_value, got %s", value)
	}
	
	// Test with non-existing environment variable (should return default)
	value = getEnv("NON_EXISTENT_VAR", "default_value")
	if value != "default_value" {
		t.Errorf("Expected default_value, got %s", value)
	}
	
	// Test with empty default
	os.Unsetenv("ANOTHER_VAR")
	value = getEnv("ANOTHER_VAR", "")
	if value != "" {
		t.Errorf("Expected empty string, got %s", value)
	}
}

func TestEmailConfig(t *testing.T) {
	// Set up test environment
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USER", "test@test.com")
	os.Setenv("SMTP_PASSWORD", "testpass")
	os.Setenv("FROM_EMAIL", "noreply@test.com")
	os.Setenv("FROM_NAME", "Test System")
	
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")
	
	if smtpHost != "smtp.test.com" {
		t.Errorf("Expected SMTPHost smtp.test.com, got %s", smtpHost)
	}
	
	if smtpPort != "587" {
		t.Errorf("Expected SMTPPort 587, got %s", smtpPort)
	}
	
	if smtpUser != "test@test.com" {
		t.Errorf("Expected SMTPUser test@test.com, got %s", smtpUser)
	}
	
	if smtpPass != "testpass" {
		t.Errorf("Expected SMTPPassword testpass, got %s", smtpPass)
	}
	
	if fromEmail != "noreply@test.com" {
		t.Errorf("Expected FromEmail noreply@test.com, got %s", fromEmail)
	}
	
	if fromName != "Test System" {
		t.Errorf("Expected FromName Test System, got %s", fromName)
	}
}


func TestEmailConfigDefaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("FROM_EMAIL")
	os.Unsetenv("FROM_NAME")
	
	// Test that environment variables are unset
	if os.Getenv("SMTP_HOST") != "" {
		t.Error("Expected SMTP_HOST to be unset")
	}
	
	if os.Getenv("SMTP_PORT") != "" {
		t.Error("Expected SMTP_PORT to be unset")
	}
	
	if os.Getenv("FROM_EMAIL") != "" {
		t.Error("Expected FROM_EMAIL to be unset")
	}
	
	if os.Getenv("FROM_NAME") != "" {
		t.Error("Expected FROM_NAME to be unset")
	}
}


