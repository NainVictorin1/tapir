package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NainVictorin1/homework2/Internal/validator"
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

// ValidateTodo ensures the fields in a to-do item are correctly formatted.
func ValidateTodo(v *validator.Validator, todo *Todo) {
	v.Check(validator.NotBlank(todo.Title), "title", "must be provided")
	v.Check(validator.MaxLength(todo.Title, 100), "title", "must not exceed 100 characters")
	v.Check(validator.NotBlank(todo.Description), "description", "must be provided")
	v.Check(validator.MaxLength(todo.Description, 1000), "description", "must not exceed 1000 characters")
	v.Check(todo.Status == "pending" || todo.Status == "completed", "status", "must be 'pending' or 'completed'")
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
func (m *TodoModel) Get(id int) (*Todo, error) {
	query := `
		SELECT id, title, description, status, created_at
		FROM todos
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	todo := &Todo{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return todo, nil
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

// Delete removes a to-do entry by ID.
func (m *TodoModel) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
