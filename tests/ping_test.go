package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func TestPing(t *testing.T) {
	AUTH_TOKEN := os.Getenv("AUTH_TOKEN_TEST")
	requestBody, err := json.Marshal(map[string]string{
		"url": "https://example.com",
	})
	if err != nil {
		t.Fatal("Failed to marshal JSON:", err)
	}
	//
	// Case 1: Create a new ping
	//
	req, err := http.NewRequest(http.MethodPost, BACKEND_URL+"/api/ping/create", bytes.NewBuffer(requestBody))
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", resp.StatusCode)
	}
	//
	// Case 2 - Create a ping with the same URL
	//
	req2, err := http.NewRequest(http.MethodPost, BACKEND_URL+"/api/ping/create", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal("Failed to make POST request:", err)
	}
	req2.Header.Set("Authorization", AUTH_TOKEN)
	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		t.Fatal("Failed to make Post request:", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, but got %d", resp2.StatusCode)
	}

	//
	// Case 3 - Delete the ping
	//
	req3, err := http.NewRequest(http.MethodDelete, BACKEND_URL+"/api/ping/delete", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal("Failed to make DELETE request:", err)
	}
	req3.Header.Set("Authorization", AUTH_TOKEN)
	client3 := &http.Client{}
	resp3, err := client3.Do(req3)
	if err != nil {
		t.Fatal("Failed to make DELETE request:", err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", resp3.StatusCode)
	}

	//
	// Case 4 - Delete the ping that does not exist
	//
	req4, err := http.NewRequest(http.MethodDelete, BACKEND_URL+"/api/ping/delete", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal("Failed to make DELETE request:", err)
	}
	req4.Header.Set("Authorization", AUTH_TOKEN)
	client4 := &http.Client{}
	resp4, err := client4.Do(req4)
	if err != nil {
		t.Fatal("Failed to make DELETE request:", err)
	}
	defer resp4.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", resp4.StatusCode)
	}
	//
	// Case 5 - Create a ping with invalid URL
	//
	requestBody2, err := json.Marshal(map[string]string{
		"url": "hs://example.com",
	})
	if err != nil {
		t.Fatal("Failed to marshal JSON:", err)
	}
	req5, err := http.NewRequest(http.MethodPost, BACKEND_URL+"/api/ping/create", bytes.NewBuffer(requestBody2))
	if err != nil {
		t.Fatal("Failed to make DELETE request:", err)
	}
	req5.Header.Set("Authorization", AUTH_TOKEN)
	client5 := &http.Client{}
	resp5, err := client5.Do(req5)
	if err != nil {
		t.Fatal("Failed to make DELETE request:", err)
	}
	defer resp5.Body.Close()
	if resp5.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, but got %d", resp5.StatusCode)
	}
}
