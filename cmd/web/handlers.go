package main

import (
	"log"
	"net/http"
)

// Home Page Handler
func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.tmpl", nil)
}

// submitFeedbackHandler handles feedback submission
func submitFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Retrieve form values
		fullname := r.FormValue("fullname")
		subject := r.FormValue("subject")
		message := r.FormValue("message")
		email := r.FormValue("email")

		// Check if required fields are filled
		if fullname == "" || message == "" {
			renderTemplate(w, "feedback.tmpl", struct {
				SuccessMessage string
				ErrorMessage   string
			}{
				ErrorMessage: "Full name and message are required.",
			})
			return
		}

		// Set default values for optional fields
		if subject == "" {
			subject = "No subject"
		}
		if email == "" {
			email = "No email provided"
		}

		// Insert feedback into the database
		_, err := database.Exec(
			"INSERT INTO feedback (fullname, subject, message, email, created_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)",
			fullname, subject, message, email,
		)
		if err != nil {
			log.Println("Error inserting feedback:", err)
			renderTemplate(w, "feedback.tmpl", struct {
				SuccessMessage string
				ErrorMessage   string
			}{
				ErrorMessage: "Failed to submit feedback. Please try again later.",
			})
			return
		}

		// Log success and render success message
		log.Println("Feedback submitted successfully")
		renderTemplate(w, "feedback.tmpl", struct {
			SuccessMessage string
			ErrorMessage   string
		}{
			SuccessMessage: "Thank you for your feedback!",
		})
		return
	}
	// Render the form for GET requests
	renderTemplate(w, "feedback.tmpl", struct {
		SuccessMessage string
		ErrorMessage   string
	}{})
}

// viewFeedbacksHandler renders a list of feedback entries
func viewFeedbacksHandler(w http.ResponseWriter, r *http.Request) {
	// Query the database for feedback entries
	rows, err := database.Query("SELECT id, fullname, subject, message, email, created_at FROM feedback ORDER BY created_at DESC")
	if err != nil {
		log.Println("Error fetching feedback:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var feedbackEntries []Feedback
	for rows.Next() {
		var f Feedback
		if err := rows.Scan(&f.ID, &f.Fullname, &f.Subject, &f.Message, &f.Email, &f.CreatedAt); err != nil {
			log.Println("Error scanning feedback row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set default values for missing fields
		if f.Subject == "" {
			f.Subject = "No subject"
		}
		if f.Email == "" {
			f.Email = "No email provided"
		}

		feedbackEntries = append(feedbackEntries, f)
	}

	// Log if no feedback was found
	if len(feedbackEntries) == 0 {
		log.Println("No feedback found")
	}

	// Render the view_feedback template with the feedback data
	renderTemplate(w, "view_feedback.tmpl", ViewFeedbackData{FeedbackEntries: feedbackEntries})
}

// submitTodoHandler handles todo item submission
func submitTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		task := r.FormValue("task")         // Use 'task' instead of 'title'
		deadline := r.FormValue("deadline") // Use 'deadline' instead of 'status'

		// Ensure task and deadline are provided
		if task == "" {
			http.Error(w, "Task cannot be empty.", http.StatusBadRequest)
			return
		}
		if deadline == "" {
			http.Error(w, "Deadline cannot be empty.", http.StatusBadRequest)
			return
		}

		// Log the values for debugging
		log.Printf("Inserting Todo - Task: %s, Deadline: %s", task, deadline)

		// Insert the Todo item into the database
		_, err := database.Exec(
			"INSERT INTO todo (task, deadline, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)",
			task, deadline,
		)
		if err != nil {
			log.Printf("Error executing query: %v", err) // Log the error
			renderTemplate(w, "add_todo.tmpl", struct{ ErrorMessage string }{"Failed to add to-do."})
			return
		}

		renderTemplate(w, "add_todo.tmpl", struct{ SuccessMessage string }{"To-do added successfully!"})
		return
	}
	renderTemplate(w, "add_todo.tmpl", nil)
}

// viewTodosHandler renders a list of todo entries
func viewTodosHandler(w http.ResponseWriter, r *http.Request) {
	// Query the todo table to fetch all required fields: task, deadline, created_at
	rows, err := database.Query("SELECT id, task, deadline, created_at FROM todo ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todoEntries []Todo
	for rows.Next() {
		var t Todo
		// Scanning the correct fields into the Todo struct
		if err := rows.Scan(&t.ID, &t.Task, &t.Deadline, &t.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		todoEntries = append(todoEntries, t)
	}

	renderTemplate(w, "view_todo.tmpl", ViewTodoData{TodoEntries: todoEntries})
}

// submitJournalHandler handles journal subject submission
func submitJournalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		subject := r.FormValue("subject")
		fullname := r.FormValue("fullname")

		if subject == "" {
			renderTemplate(w, "submit_journal.tmpl", SubmitJournalData{ErrorMessage: "Subject is required."})
			return
		}
		if fullname == "" {
			fullname = "No name provided"
		}

		_, err := database.Exec(
			"INSERT INTO journal (subject, fullname, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)",
			subject, fullname,
		)
		if err != nil {
			renderTemplate(w, "submit_journal.tmpl", SubmitJournalData{ErrorMessage: "Failed to submit journal entry."})
			return
		}

		renderTemplate(w, "submit_journal.tmpl", SubmitJournalData{SuccessMessage: "Journal entry submitted successfully!"})
		return
	}
	renderTemplate(w, "submit_journal.tmpl", nil)
}

// listJournalsHandler renders a list of journal entries
func viewJournalsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT id, subject, fullname, created_at FROM journal ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var journalEntries []Journal
	for rows.Next() {
		var j Journal
		err := rows.Scan(&j.ID, &j.Subject, &j.Fullname, &j.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Optional: fallback values

		if j.Subject == "" {
			j.Subject = "No subject"
		}

		journalEntries = append(journalEntries, j)
	}

	renderTemplate(w, "view_journal.tmpl", ViewJournalData{JournalEntries: journalEntries})
}
