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

func TestPing(t *testing.T) {
	resp, err := http.Post(BACKEND_URL+"/ping/create", "application/json", nil)
	if err != nil {
		t.Fatal("Failed to make POST request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, but got %d", resp.StatusCode)
	}
}
