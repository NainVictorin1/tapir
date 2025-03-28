package main

import (
	"net/http"

	"github.com/NainVictorin1/homework2/Internal/data"
	"github.com/NainVictorin1/homework2/Internal/validator"
)

// Home Page Handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Welcome"
	data.HeaderText = "We are here to help"

	err := app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// feedbackFormHandler renders the feedback submission form
func (app *application) feedbackHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Submit Feedback"

	err := app.render(w, http.StatusOK, "feedback.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render feedback form", "template", "feedback.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// submitFeedbackHandler handles feedback submission
func (app *application) submitFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	feedback := &data.Feedback{
		Fullname: r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Subject:  r.PostForm.Get("subject"),
		Message:  r.PostForm.Get("message"),
	}

	v := validator.NewValidator()
	data.ValidateFeedback(v, feedback)

	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Submit Feedback"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"name":    feedback.Fullname,
			"email":   feedback.Email,
			"subject": feedback.Subject,
			"message": feedback.Message,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "feedback.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render feedback form", "template", "feedback.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.feedback.Insert(feedback)
	if err != nil {
		app.logger.Error("failed to insert feedback entry", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/feedback/success", http.StatusSeeOther)
}

// viewFeedbacksHandler renders a list of feedback entries
func (app *application) viewFeedbacksHandler(w http.ResponseWriter, r *http.Request) {
	feedbacks, err := app.feedback.GetAll()
	if err != nil {
		app.logger.Error("failed to fetch feedback entries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "View Feedbacks"
	data.Feedbacks = feedbacks

	err = app.render(w, http.StatusOK, "view_feedback.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render feedbacks page", "template", "view_feedback.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Todo Handlers

// todoFormHandler renders the todo submission form
func (app *application) todoHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Add Todo"

	err := app.render(w, http.StatusOK, "add_todo.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render todo form", "template", "add_todo.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// submitTodoHandler handles todo item submission
func (app *application) submitTodoHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	todo := &data.Todo{
		Title:       r.PostForm.Get("titulo"),
		Description: r.PostForm.Get("description"),
	}

	v := validator.NewValidator()
	data.ValidateTodo(v, todo)

	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Add Todo"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"titulo":      todo.Title,
			"description": todo.Description,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "add_todo.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render todo form", "template", "add_todo.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.todos.Insert(todo)
	if err != nil {
		app.logger.Error("failed to insert todo item", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/todo/view", http.StatusSeeOther)
}

// viewTodosHandler renders a list of todo entries
func (app *application) viewTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := app.todos.GetAll()
	if err != nil {
		app.logger.Error("failed to fetch todo entries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "View Todos"
	data.Todos = todos

	err = app.render(w, http.StatusOK, "view_todo.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render todos page", "template", "view_todo.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) JournalHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Submit a Journal Entry"

	err := app.render(w, http.StatusOK, "submit_journal.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render journal form", "template", "submit_journal.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// submitJournalHandler handles journal entry submission
func (app *application) submitJournalHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	journal := &data.Journal{
		Title: r.PostForm.Get("title"),
		Entry: r.PostForm.Get("entry"),
	}

	v := validator.NewValidator()
	data.ValidateJournal(v, journal)

	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Submit a Journal Entry"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"title": journal.Title,
			"entry": journal.Entry,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "submit_journal.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render journal form", "template", "submit_journal.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.journals.Insert(journal) // Ensure this matches your struct name
	if err != nil {
		app.logger.Error("failed to insert journal entry", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/journals", http.StatusSeeOther)
}

// listJournalsHandler renders a list of journal entries
func (app *application) viewJournalsHandler(w http.ResponseWriter, r *http.Request) {
	journals, err := app.journals.GetAll()
	if err != nil {
		app.logger.Error("failed to fetch journal entries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "ViewJournal "
	data.Journals = journals

	err = app.render(w, http.StatusOK, "view_journal.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render journals page", "template", "view_journal.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
