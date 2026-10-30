package memory

import (
	"context"
	"sort"
	"sync"

	"api/internal/app/tasks/domain"
)

type Repository struct {
	mu    sync.RWMutex
	tasks map[string]domain.Task
}

var _ domain.Repository = (*Repository)(nil)

func New() *Repository { return &Repository{tasks: make(map[string]domain.Task)} }
func (r *Repository) Create(ctx context.Context, task domain.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task

	return nil
}
func (r *Repository) Get(ctx context.Context, id string) (domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return domain.Task{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}

	return task, nil
}
func (r *Repository) List(ctx context.Context) ([]domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].CreatedAt.Equal(tasks[j].CreatedAt) {
			return tasks[i].ID < tasks[j].ID
		}

		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	return tasks, nil
}
