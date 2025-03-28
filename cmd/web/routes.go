package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Serve static files from the ./ui/static/ directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Home Route
	mux.HandleFunc("/", app.homeHandler) // Corrected route for the home handler

	// Feedback Routes
	mux.HandleFunc("/feedback/submit", app.submitFeedbackHandler)
	mux.HandleFunc("/feedback/view", app.viewFeedbacksHandler) // Displays all feedback entries
	mux.HandleFunc("/feedback/create", app.submitFeedbackHandler)

	// Todo Routes
	mux.HandleFunc("/todo/add", app.addTodoHandler)    // Handles displaying the to-do form
	mux.HandleFunc("/todo/create", app.addTodoHandler) // Handles to-do item submission
	mux.HandleFunc("/todo/view", app.viewTodosHandler) // Displays all to-do items

	// Journal Routes
	mux.HandleFunc("/journal/add", app.submitJournalHandler)    // Handles displaying the journal form
	mux.HandleFunc("/journal/create", app.submitJournalHandler) // Handles journal submission
	mux.HandleFunc("/journal/view", app.viewJournalsHandler)    // Displays all journal entries

	// Return the middleware-wrapped mux
	return app.loggingMiddleware(mux)
}
