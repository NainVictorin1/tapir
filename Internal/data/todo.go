package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Todo represents a to-do item in the database.
type Todo struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}

// TodoModel handles database operations for to-do entries.
type TodoModel struct {
	DB *sql.DB
}

// Insert adds a new to-do entry to the database.
func (m *TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todos (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		todo.Title,
		todo.Description,
		todo.Status,
	).Scan(&todo.ID, &todo.CreatedAt)
}

// GetAll retrieves all to-do items sorted by ID (newest first).
func (m *TodoModel) GetAll() ([]*Todo, error) {
	query := `
		SELECT id, title, description, status, created_at
		FROM todos
		ORDER BY id DESC`

	fmt.Println("Fetching to-do data...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Status,
			&todo.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	fmt.Println("Successfully retrieved to-do data.")
	return todos, nil
}
