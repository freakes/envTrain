package postgres

import (
	"database/sql"
	"relay/internal/app/tasks/domain"
)

type Repository struct{ db *sql.DB }

func New(db *sql.DB) *Repository { return &Repository{db: db} }

var _ domain.Repository = (*Repository)(nil)

type batch struct {
	tx     *sql.Tx
	events []domain.Event
}
