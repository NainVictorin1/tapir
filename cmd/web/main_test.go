package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type application struct {
	logger *slog.Logger
}

func (app *application) viewTodosHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation of viewTodosHandler
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo List"))
}

func (app *application) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation of addTodoHandler
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo added successfully!"))
}

// TestAddTodoHandler tests the addTodoHandler function for adding a to-do item.
func TestAddTodoHandler(t *testing.T) {
	// Set up logger and application instance
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	// Test adding a to-do with task, name, and email
	req := httptest.NewRequest("POST", "/todo", strings.NewReader("task=Finish writing Go tests&name=John Doe&email=johndoe@example.com"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.addTodoHandler)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
	}

	// Check if the response contains a success message
	expected := "To-do added successfully!"
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

// TestViewTodosHandler tests the viewTodosHandler function for viewing to-do items.
func TestViewTodosHandler(t *testing.T) {
	// Set up logger and application instance
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	// Simulate a GET request to view the to-do list
	req := httptest.NewRequest("GET", "/todos", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.viewTodosHandler)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
	}

	// Check if the response contains the "To-do List"
	expected := "To-do List"
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

// TestAddTodoHandler_InvalidTask tests the scenario where the task field is empty.
func TestAddTodoHandler_InvalidTask(t *testing.T) {
	// Set up logger and application instance
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	// Simulate a POST request with an empty task field
	req := httptest.NewRequest("POST", "/todo", strings.NewReader("task=&name=John Doe&email=johndoe@example.com"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.addTodoHandler)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusBadRequest {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusBadRequest)
	}

	// Check for the error message in the response body
	expected := "Task cannot be empty."
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

// TestAddTodoHandler_InvalidFormData tests the scenario where name or email is missing.
func TestAddTodoHandler_InvalidFormData(t *testing.T) {
	// Set up logger and application instance
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	// Simulate a POST request with missing name and email
	req := httptest.NewRequest("POST", "/todo", strings.NewReader("task=Test task"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.addTodoHandler)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusBadRequest {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusBadRequest)
	}

	// Check for the error message in the response body
	expected := "Name and Email are required."
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

// TestDatabaseConnection tests the database connection.
func TestDatabaseConnection(t *testing.T) {
	// Initialize the database connection (initDB should be a valid function)
	initdatabase()

	// Check if the database is closed after the defer
	defer database.Close()
	if database == nil {
		t.Fatalf("Database connection is nil")
	}
}
