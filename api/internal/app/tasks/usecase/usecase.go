package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"api/internal/app/tasks/domain"
)

type usecase struct{ repo domain.Repository }

func New(repo domain.Repository) domain.Usecase { return &usecase{repo: repo} }
func (u *usecase) Create(ctx context.Context, req domain.CreateRequest) (domain.Task, error) {
	if err := req.Validate(); err != nil {
		return domain.Task{}, err
	}

	task := domain.Task{ID: rand.Text(), Title: strings.TrimSpace(req.Title), CreatedAt: time.Now().UTC().Truncate(time.Microsecond)}
	if err := u.repo.Create(ctx, task); err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
func (u *usecase) Get(ctx context.Context, id string) (domain.Task, error) {
	task, err := u.repo.Get(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	return task, nil
}
func (u *usecase) List(ctx context.Context) ([]domain.Task, error) {
	tasks, err := u.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, nil
}
