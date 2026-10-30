package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"api/internal/app/tasks/domain"
	"api/internal/app/tasks/repository/memory"
	"api/internal/app/tasks/usecase"
)

func TestCreateValidationAndStorage(t *testing.T) {
	for _, tt := range []struct {
		name, title string
		valid       bool
	}{
		{"empty", "", false}, {"blank", " \n\t", false}, {"too long", strings.Repeat("я", 201), false},
		{"unicode boundary", strings.Repeat("я", 200), true}, {"trim", "  task  ", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.New()
			uc := usecase.New(repo)
			task, err := uc.Create(context.Background(), domain.CreateRequest{Title: tt.title})
			if !tt.valid {
				if !errors.Is(err, domain.ErrInvalidTitle) {
					t.Fatalf("want validation error, got %v", err)
				}
				all, _ := repo.List(context.Background())
				if len(all) != 0 {
					t.Fatal("invalid task persisted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if task.ID == "" || task.CreatedAt.IsZero() || task.Title != strings.TrimSpace(tt.title) {
				t.Fatalf("unexpected task: %+v", task)
			}
			stored, err := uc.Get(context.Background(), task.ID)
			if err != nil || stored != task {
				t.Fatalf("stored task mismatch: %+v %v", stored, err)
			}
		})
	}
}
func TestWrappedErrors(t *testing.T) {
	uc := usecase.New(memory.New())
	if _, err := uc.Get(context.Background(), "missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("got %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := uc.Create(ctx, domain.CreateRequest{Title: "task"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
}
