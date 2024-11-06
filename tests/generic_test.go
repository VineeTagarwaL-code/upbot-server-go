package tests

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	resp, err := http.Get(BACKEND_URL + "/api/health")
	if err != nil {
		t.Fatal("Failed to make GET request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, but got %d", resp.StatusCode)
	}
}
