package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Set up application with logger
	app := &application{
		logger: logger,
	}

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/", nil)

	// Create a new ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Use the correct handler function (app.homeHandler in your code)
	handler := http.HandlerFunc(app.homeHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check if the status code is 200 OK
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
	}

	// Define the expected output based on what the home template renders
	expected := "<html><body><h1>Welcome</h1><p>We are here to help</p></body></html>" // Adjust according to actual rendered output
	got := rr.Body.String()

	// Compare the expected and actual response
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}
