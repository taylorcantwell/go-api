package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestPostMessage(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatalf("Failed to open in-memory sqlite database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	setupHandlers(mux, db)

	server := httptest.NewServer(mux)
	defer server.Close()

	bodyBytes, _ := json.Marshal(RequestBody{Message: "Hello, Test!"})
	request, err := http.NewRequest("POST", server.URL+"/message", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, response.StatusCode)
	}

	var respBody Response
	json.NewDecoder(response.Body).Decode(&respBody)
	if respBody.Message != "Hello, Test!" {
		t.Errorf("Expected message 'Hello, Test!', got '%s'", respBody.Message)
	}
}
