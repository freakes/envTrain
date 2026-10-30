package postgres

import (
	"context"
	"database/sql"

	"outbox-consumer/internal/app/tasks/domain"
)

type Repository struct{ db *sql.DB }

func New(db *sql.DB) *Repository { return &Repository{db: db} }

var _ domain.Repository = (*Repository)(nil)
