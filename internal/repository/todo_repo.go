package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ujjvalkumar/todo-api/internal/domain"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	query := `INSERT INTO todos (id, user_id, title, description, status, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		todo.ID.String(),
		todo.UserID.String(),
		todo.Title,
		todo.Description,
		todo.Status,
		todo.CreatedAt,
		todo.UpdatedAt,
	)
	return err
}

func (r *TodoRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Todo, error) {
	query := `SELECT id, user_id, title, description, status, created_at, updated_at
			  FROM todos WHERE id = ? AND user_id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String(), userID.String())

	var t domain.Todo
	var idStr, userIDStr string
	err := row.Scan(&idStr, &userIDStr, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	t.ID, _ = uuid.Parse(idStr)
	t.UserID, _ = uuid.Parse(userIDStr)
	return &t, nil
}

func (r *TodoRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error) {
	query := `SELECT id, user_id, title, description, status, created_at, updated_at
			  FROM todos WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []domain.Todo
	for rows.Next() {
		var t domain.Todo
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		t.ID, _ = uuid.Parse(idStr)
		t.UserID, _ = uuid.Parse(userIDStr)
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (r *TodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `UPDATE todos SET title = ?, description = ?, status = ?, updated_at = ?
			  WHERE id = ? AND user_id = ?`
	result, err := r.db.ExecContext(ctx, query,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.UpdatedAt,
		todo.ID.String(),
		todo.UserID.String(),
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM todos WHERE id = ? AND user_id = ?`
	result, err := r.db.ExecContext(ctx, query, id.String(), userID.String())
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
