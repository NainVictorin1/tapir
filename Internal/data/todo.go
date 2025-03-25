package data

import (
	"context"
	"database/sql"
	"time"
)

// Todo represents a todo item in the database
type Todo struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Deadline  time.Time `json:"deadline"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoModel handles database operations for todo items
type TodoModel struct {
	DB *sql.DB
}

// Insert adds a new todo item to the database
func (m *TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todos (task, deadline, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, todo.Task, todo.Deadline).Scan(&todo.ID, &todo.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a specific todo item by ID
func (m *TodoModel) Get(id int) (*Todo, error) {
	query := `
		SELECT id, task, deadline, created_at
		FROM todos
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	todo := &Todo{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Task,
		&todo.Deadline,
		&todo.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// GetAll retrieves all todo items sorted by deadline (soonest first)
func (m *TodoModel) GetAll() ([]*Todo, error) {
	query := `
		SELECT id, task, deadline, created_at
		FROM todos
		ORDER BY deadline ASC
	`

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
			&todo.Task,
			&todo.Deadline,
			&todo.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// Delete removes a todo item from the database by ID
func (m *TodoModel) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
