package utils

import "os"

// GetAuthServiceURL returns the auth service URL
func GetAuthServiceURL() string {
	url := os.Getenv("AUTH_SERVICE_URL")
	if url == "" {
		return "http://localhost:8001"
	}
	return url
}

// GetBuildingServiceURL returns the building service URL
func GetBuildingServiceURL() string {
	url := os.Getenv("BUILDING_SERVICE_URL")
	if url == "" {
		return "http://localhost:8002"
	}
	return url
}
