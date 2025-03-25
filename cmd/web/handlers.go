package main

import (
	"net/http"
	"time"

	"github.com/NainVictorin1/homework2/Internal/data"
	"github.com/NainVictorin1/homework2/Internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Welcome"
	data.HeaderText = "We are here to help"
	err := app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) createFeedback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	subject := r.PostForm.Get("subject")
	message := r.PostForm.Get("message")

	feedback := &data.Feedback{
		Fullname: name,
		Email:    email,
		Subject:  subject,
		Message:  message,
	}

	// Validate data
	v := validator.NewValidator()
	data.ValidateFeedback(v, feedback)
	// Check for validation errors
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Welcome"
		data.HeaderText = "We are here to help"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"name":    name,
			"email":   email,
			"subject": subject,
			"message": message,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "home.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.feedback.Insert(feedback)
	if err != nil {
		app.logger.Error("failed to insert feedback", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/feedback/success", http.StatusSeeOther)
}

func (app *application) feedbackSuccess(w http.ResponseWriter, r *http.Request) {
	templateData := NewTemplateData()
	templateData.Title = "Feedback Submitted"
	templateData.HeaderText = "Thank You for Your Feedback!"

	err := app.render(w, http.StatusOK, "feedback_success.tmpl", templateData)
	if err != nil {
		app.logger.Error("failed to render feedback success page", "template", "feedback_success.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) addTodoForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Add Todo"

	err := app.render(w, http.StatusOK, "add_todo.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render add_todo page", "template", "add_todo.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Process form submission for adding a new Todo
func (app *application) createTodo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse todo form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	task := r.PostForm.Get("task")
	deadlineStr := r.PostForm.Get("deadline")

	deadline, err := time.Parse("2006-01-02", deadlineStr) // ✅ Use proper date format
	if err != nil {
		app.logger.Error("invalid deadline format", "error", err)
		http.Error(w, "Invalid date format. Please use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	todo := &data.Todo{
		Task:     task,
		Deadline: deadline,
	}

	// Validate data
	v := validator.NewValidator()
	data.ValidateTodo(v, todo)

	if !v.ValidData() {
		templateData := NewTemplateData()
		templateData.Title = "Add Todo"
		templateData.FormErrors = v.Errors
		templateData.FormData = map[string]string{
			"task":     task,
			"deadline": deadlineStr,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "add_todo.tmpl", templateData)
		if err != nil {
			app.logger.Error("failed to render add_todo page", "template", "add_todo.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.todos.Insert(todo)
	if err != nil {
		app.logger.Error("failed to insert todo", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/todo/view", http.StatusSeeOther)
}

// View all Todos
func (app *application) viewTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := app.todos.GetAll()
	if err != nil {
		app.logger.Error("failed to fetch todos", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "View Todos"
	data.Todos = todos

	err = app.render(w, http.StatusOK, "view_todo.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render view_todo page", "template", "view_todo.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Show the form to add a Journal entry
func (app *application) addJournalForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Add Journal Entry"

	err := app.render(w, http.StatusOK, "view_journal.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render view_journal page", "template", "view_journal.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Process form submission for adding a new Journal entry
func (app *application) createJournal(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse journal form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	entry := r.PostForm.Get("entry")

	journal := &data.Journal{
		Title: title,
		Entry: entry,
	}

	// Validate data
	v := validator.NewValidator()
	data.ValidateJournal(v, journal)

	if !v.ValidData() {
		templateData := NewTemplateData()
		templateData.Title = "Add Journal Entry"
		templateData.FormErrors = v.Errors
		templateData.FormData = map[string]string{
			"title": title,
			"entry": entry,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "view_journal.tmpl", templateData)
		if err != nil {
			app.logger.Error("failed to render view_journal page", "template", "view_journal.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.journals.Insert(journal)
	if err != nil {
		app.logger.Error("failed to insert journal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/journal/view", http.StatusSeeOther)
}
