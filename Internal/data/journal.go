package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NainVictorin1/homework2/Internal/validator"
	_ "github.com/lib/pq"
)

// Journal represents a journal entry.
type Journal struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Entry     string    `json:"entry"`
}

// ValidateJournal ensures journal fields meet required constraints.
func ValidateJournal(v *validator.Validator, journal *Journal) {
	v.Check(validator.NotBlank(journal.Title), "title", "must be provided")
	v.Check(validator.MaxLength(journal.Title, 100), "title", "must not exceed 100 characters")
	v.Check(validator.NotBlank(journal.Entry), "entry", "must be provided")
	v.Check(validator.MaxLength(journal.Entry, 1000), "entry", "must not exceed 1000 characters")
}

// JournalModel handles database operations for journal entries.
type JournalModel struct {
	DB *sql.DB
}

// Insert adds a new journal entry to the database.
func (m *JournalModel) Insert(journal *Journal) error {
	query := `
		INSERT INTO journals (title, entry)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, journal.Title, journal.Entry).Scan(&journal.ID, &journal.CreatedAt)
}

// Get retrieves a specific journal entry by ID.
func (m *JournalModel) Get(id int64) (*Journal, error) {
	query := `
		SELECT id, title, entry, created_at
		FROM journals
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	journal := &Journal{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&journal.ID,
		&journal.Title,
		&journal.Entry,
		&journal.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return journal, nil
}

// GetAll retrieves all journal entries sorted by ID (newest first).
func (m *JournalModel) GetAll() ([]*Journal, error) {
	query := `
		SELECT id, title, entry, created_at
		FROM journals
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
			&journal.Title,
			&journal.Entry,
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

// Delete removes a journal entry by ID.
func (m *JournalModel) Delete(id int64) error {
	query := `DELETE FROM journals WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
