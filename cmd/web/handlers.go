package main

import (
	"net/http"
	"time"

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
func (app *application) addFeedbackForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Add Feedback"

	err := app.render(w, http.StatusOK, "feedback.tmpl", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ✅ Feedback Handlers
func (app *application) createFeedback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
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
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.feedback.Insert(feedback)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/feedback/success", http.StatusSeeOther)
}

func (app *application) viewFeedbacks(w http.ResponseWriter, r *http.Request) {
	feedbacks, err := app.feedback.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "View Feedbacks"
	data.Feedbacks = feedbacks

	err = app.render(w, http.StatusOK, "view_feedback.tmpl", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Todo Handlers
func (app *application) addTodoForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Add Todo"

	err := app.render(w, http.StatusOK, "add_todo.tmpl", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) createTodo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	task := r.PostForm.Get("task")
	deadlineStr := r.PostForm.Get("deadline")

	deadline, err := time.Parse("2006-01-02", deadlineStr)
	if err != nil {
		http.Error(w, "Invalid date format. Please use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	todo := &data.Todo{
		Task:     task,
		Deadline: deadline,
	}

	v := validator.NewValidator()
	data.ValidateTodo(v, todo)

	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Add Todo"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"task":     task,
			"deadline": deadlineStr,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "add_todo.tmpl", data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.todos.Insert(todo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/todo/view", http.StatusSeeOther)
}

func (app *application) viewTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := app.todos.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "View Todos"
	data.Todos = todos

	err = app.render(w, http.StatusOK, "view_todo.tmpl", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ✅ Journal Handlers
func (app *application) addJournalForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Add Journal Entry"

	err := app.render(w, http.StatusOK, "submit_journal.tmpl", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) createJournal(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
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
		data.Title = "Add Journal Entry"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"title": journal.Title,
			"entry": journal.Entry,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "submit_journal.tmpl", data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.journals.Insert(journal)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/journal/view", http.StatusSeeOther)
}

func (app *application) viewJournals(w http.ResponseWriter, r *http.Request) {
	journals, err := app.journals.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "View Journals"
	data.Journals = journals

	err = app.render(w, http.StatusOK, "view_journal.tmpl", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
