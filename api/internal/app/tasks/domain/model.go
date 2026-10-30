package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrNotFound     = errors.New("task not found")
	ErrInvalidTitle = errors.New("title must contain 1 to 200 characters")
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequest struct {
	Title string `json:"title"`
}

func (r CreateRequest) Validate() error {
	n := utf8.RuneCountInString(strings.TrimSpace(r.Title))
	if n == 0 || n > 200 {
		return ErrInvalidTitle
	}

	return nil
}
