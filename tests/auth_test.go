package tests

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var BACKEND_URL string

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	BACKEND_URL = os.Getenv("BACKEND_URL")
	if BACKEND_URL == "" {
		log.Fatal("BACKEND_URL is not set in .env file")
	}

	code := m.Run()
	os.Exit(code)
}

func TestAuth(t *testing.T) {
	AUTH_TOKEN := "Bearer " + "ADJLFLASJFDLASJFLJL"
	req, err := http.NewRequest(http.MethodGet, BACKEND_URL+"/api/auth/google", nil)
	if err != nil {
		t.Fatal("Failed to make POST request:", err)
	}
	req.Header.Set("Authorization", AUTH_TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to make GET request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code 401, but got %d", resp.StatusCode)
	}

}
