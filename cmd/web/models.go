package main

import (
	"database/sql"
	"time"
)

// Journal Struct represents a journal entry
type Journal struct {
	ID        int            // Unique identifier for the journal entry
	Title     string         // Title of the journal entry
	Subject   string         // Content of the journal entry
	CreatedAt time.Time      // Timestamp of when the journal entry was created
	Fullname  sql.NullString // Name of the user who submitted the journal entry
}

// Feedback Struct represents a feedback entry
type Feedback struct {
	ID        int            // Unique identifier for the feedback entry
	Subject   string         // Name of the person who submitted the feedback
	Message   string         // Feedback message
	Fullname  sql.NullString // Email of the person who submitted the feedback
	CreatedAt time.Time
	Email     string // Timestamp of when the feedback entry was created

}

// Todo Struct represents a to-do item
type Todo struct {
	ID        int
	Task      string
	Deadline  string
	CreatedAt string
}

// Data structures to pass journal entries to templates
type ViewJournalData struct {
	JournalEntries []Journal // List of journal entries to be displayed
}

// Data structures to pass feedback entries to templates
type ViewFeedbackData struct {
	FeedbackEntries []Feedback // List of feedback entries to be displayed
}

// Data structures to pass to-do entries to templates
type ViewTodoData struct {
	TodoEntries []Todo // List of to-do entries to be displayed
}

// SubmitJournalData holds the success or error message data to be passed to the template
type SubmitJournalData struct {
	SuccessMessage string // Message to display when a journal entry is successfully submitted
	ErrorMessage   string // Message to display when there is an error submitting the journal entry
}
