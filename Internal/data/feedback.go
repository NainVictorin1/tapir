package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

type FeedbackModel struct {
	DB *sql.DB
}

func (m *FeedbackModel) Insert(feedback *Feedback) error {
	query := `
	INSERT INTO feedback (fullname, email, subject, message) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, feedback.Fullname, feedback.Email, feedback.Subject, feedback.Message).Scan(&feedback.ID, &feedback.CreatedAt)
}
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
