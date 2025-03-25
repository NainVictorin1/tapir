package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("/todo/add", app.addTodoForm)
	mux.HandleFunc("/todo/create", app.createTodo)
	mux.HandleFunc("/todo/view", app.viewTodos)

	mux.HandleFunc("/journal/add", app.addJournalForm)
	mux.HandleFunc("/journal/create", app.createJournal)

	mux.HandleFunc("/feedback/success", app.feedbackSuccess)

	return app.loggingMiddleware(mux)
}
