package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NainVictorin1/homework2/Internal/validator"
	_ "github.com/lib/pq"
)

type Feedback struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Fullname  string    `json:"fullname"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Email     string    `json:"email"`
}

func ValidateFeedback(v *validator.Validator, feedback *Feedback) {
	v.Check(validator.NotBlank(feedback.Fullname), "fullname", "must be provided")
	v.Check(validator.MaxLength(feedback.Fullname, 50), "fullname", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(feedback.Subject), "subject", "must be provided")
	v.Check(validator.MaxLength(feedback.Subject, 50), "subject", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(feedback.Email), "email", "must be provided")
	v.Check(validator.IsValidEmail(feedback.Email), "email", "invalid email address")
	v.Check(validator.MaxLength(feedback.Email, 100), "email", "must not be more than 100 bytes long")
	v.Check(validator.NotBlank(feedback.Message), "message", "must be provided")
	v.Check(validator.MaxLength(feedback.Message, 500), "message", "must not be more than 500 bytes long")
}

type FeedbackModel struct {
	DB *sql.DB
}

func (m *FeedbackModel) Insert(feedback *Feedback) error {
	query := `
		INSERT INTO feedback (fullname, subject, message, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		feedback.Fullname,
		feedback.Subject,
		feedback.Message,
		feedback.Email,
	).Scan(&feedback.ID, &feedback.CreatedAt)
}
func (m *FeedbackModel) Get(id int) (*Feedback, error) {
	query := `
		SELECT id, title, entry, created_at
		FROM feedback
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	feedback := &Feedback{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&feedback.ID,
		&feedback.Fullname,
		&feedback.Subject,
		&feedback.Message,
		&feedback.Email,
		&feedback.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return feedback, nil
}

// GetAll retrieves all journal entries sorted by date (newest first)
func (m *FeedbackModel) GetAll() ([]*Feedback, error) {
	query := `
		SELECT id, fullname, subject, message, email, created_at
		FROM feedback
		ORDER BY id DESC`

	fmt.Println("Fetching feedback data...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []*Feedback
	for rows.Next() {
		feedback := &Feedback{}
		err := rows.Scan(
			&feedback.ID,
			&feedback.Fullname,
			&feedback.Subject,
			&feedback.Message,
			&feedback.Email,
			&feedback.CreatedAt, // Ensure this is the last field
		)
		if err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	fmt.Println("Successfully retrieved feedback data.")
	return feedbacks, nil
}

// Delete removes a journal entry from the database by ID
func (m *FeedbackModel) Delete(id int) error {
	query := `DELETE FROM feedback WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
