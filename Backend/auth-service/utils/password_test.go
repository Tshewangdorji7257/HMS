package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if hash == "" {
		t.Fatal("Expected hash to not be empty")
	}
	
	if hash == password {
		t.Fatal("Hash should not equal original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	// Test with correct password
	if !CheckPasswordHash(password, hash) {
		t.Error("Expected password to match hash")
	}
	
	// Test with incorrect password
	if CheckPasswordHash("wrongPassword", hash) {
		t.Error("Expected password to not match hash")
	}
	
	// Test with empty password
	if CheckPasswordHash("", hash) {
		t.Error("Expected empty password to not match hash")
	}
}

func TestHashPasswordDifferentPasswords(t *testing.T) {
	password1 := "password123"
	password2 := "password456"
	
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)
	
	if hash1 == hash2 {
		t.Error("Different passwords should have different hashes")
	}
}

func TestHashPasswordSamePasswordDifferentHashes(t *testing.T) {
	password := "testPassword"
	
	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)
	
	// Due to bcrypt salting, same password should produce different hashes
	if hash1 == hash2 {
		t.Error("Same password should produce different hashes due to salting")
	}
	
	// But both should be valid
	if !CheckPasswordHash(password, hash1) || !CheckPasswordHash(password, hash2) {
		t.Error("Both hashes should be valid for the same password")
	}
}
