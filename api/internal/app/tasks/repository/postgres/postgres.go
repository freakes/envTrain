package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"api/internal/app/tasks/domain"
)

type Repository struct{ db *sql.DB }

var _ domain.Repository = (*Repository)(nil)

func New(db *sql.DB) *Repository { return &Repository{db: db} }
func (r *Repository) Create(ctx context.Context, task domain.Task) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO tasks (id,title,created_at) VALUES ($1,$2,$3)`, task.ID, task.Title, task.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	return nil
}
func (r *Repository) Get(ctx context.Context, id string) (domain.Task, error) {
	var task domain.Task

	err := r.db.QueryRowContext(ctx, `SELECT id,title,created_at FROM tasks WHERE id=$1`, id).
		Scan(&task.ID, &task.Title, &task.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return task, domain.ErrNotFound
	}

	if err != nil {
		return task, fmt.Errorf("select task: %w", err)
	}

	return task, nil
}
func (r *Repository) List(ctx context.Context) ([]domain.Task, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,title,created_at FROM tasks ORDER BY created_at,id`)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)

	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.CreatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}
