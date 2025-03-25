package main

import (
	"net/http"

	"github.com/NainVictorin1/homework2/Internal/Data"
	"github.com/NainVictorin1/homework2/Internal/Data/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	templateData := NewTemplateData()
	templateData.Title = "Welcome"
	templateData.HeaderText = "We are here to help"
	err := app.render(w, http.StatusOK, "home.tmpl", templateData)
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

	feedback := &Data.Feedback{
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
