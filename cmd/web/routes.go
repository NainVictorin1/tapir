package main

import (
	"net/http"
)

func registerRoutes() {
	http.HandleFunc("/", loggingMiddleware(http.HandlerFunc(home)).ServeHTTP)
	http.HandleFunc("/home", loggingMiddleware(http.HandlerFunc(home)).ServeHTTP)
	http.HandleFunc("/journal/add", loggingMiddleware(http.HandlerFunc(submitJournalHandler)).ServeHTTP)
	http.HandleFunc("/journal/view", loggingMiddleware(http.HandlerFunc(viewJournalsHandler)).ServeHTTP)
	http.HandleFunc("/feedback/submit", loggingMiddleware(http.HandlerFunc(submitFeedbackHandler)).ServeHTTP)
	http.HandleFunc("/feedback/view", loggingMiddleware(http.HandlerFunc(viewFeedbacksHandler)).ServeHTTP)
	http.HandleFunc("/todo/add", loggingMiddleware(http.HandlerFunc(submitTodoHandler)).ServeHTTP)
	http.HandleFunc("/todo/view", loggingMiddleware(http.HandlerFunc(viewTodosHandler)).ServeHTTP)

	// Static file server
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
