package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Journal represents a journal entry.
type Journal struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Fullname  string    `json:"fullname"`
	Subject   string    `json:"subject"`
}

// JournalModel handles database operations for journal entries.
type JournalModel struct {
	DB *sql.DB
}

// Insert adds a new journal entry to the database.
func (m *JournalModel) Insert(journal *Journal) error {
	query := `
		INSERT INTO journal (fullname, subject)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, journal.Fullname, journal.Subject).Scan(&journal.ID, &journal.CreatedAt)
}

// GetAll retrieves all journal entries sorted by ID (newest first).
func (m *JournalModel) GetAll() ([]*Journal, error) {
	query := `
		SELECT id, fullname, subject, created_at
		FROM journal
		ORDER BY id DESC`

	fmt.Println("Fetching journal data...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []*Journal
	for rows.Next() {
		journal := &Journal{}
		err := rows.Scan(
			&journal.ID,
			&journal.Fullname,
			&journal.Subject,
			&journal.CreatedAt, // Ensure this is the last field scanned
		)
		if err != nil {
			return nil, err
		}
		journals = append(journals, journal)
	}

	fmt.Println("Successfully retrieved journal data.")
	return journals, nil
}
