package Data

import (
	"context"
	"database/sql"
	"time"
)

// Journal represents a journal entry in the database
type Journal struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Entry     string    `json:"entry"`
	CreatedAt time.Time `json:"created_at"`
}

// JournalModel handles database operations for journal entries
type JournalModel struct {
	DB *sql.DB
}

// Insert adds a new journal entry to the database
func (m *JournalModel) Insert(journal *Journal) error {
	query := `
		INSERT INTO journals (title, entry, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, journal.Title, journal.Entry).Scan(&journal.ID, &journal.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a specific journal entry by ID
func (m *JournalModel) Get(id int) (*Journal, error) {
	query := `
		SELECT id, title, entry, created_at
		FROM journals
		WHERE id = $1
	`

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

// GetAll retrieves all journal entries sorted by date (newest first)
func (m *JournalModel) GetAll() ([]*Journal, error) {
	query := `
		SELECT id, title, entry, created_at
		FROM journals
		ORDER BY created_at DESC
	`

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
			&journal.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		journals = append(journals, journal)
	}

	return journals, nil
}

// Delete removes a journal entry from the database by ID
func (m *JournalModel) Delete(id int) error {
	query := `DELETE FROM journals WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
