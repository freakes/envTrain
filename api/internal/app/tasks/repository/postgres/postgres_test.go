package postgres_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"api/internal/app/tasks/domain"
	repository "api/internal/app/tasks/repository/postgres"
	"api/internal/app/tasks/usecase"
	"api/pkg/postgres"
)

// Run against a disposable, migrated database; no production data is deleted.
func TestPostgresLifecycle(t *testing.T) {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("set TEST_POSTGRES_DSN to run PostgreSQL integration test")
	}
	ctx := context.Background()
	db, err := postgres.Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	uc := usecase.New(repository.New(db))
	task, err := uc.Create(ctx, domain.CreateRequest{Title: "integration test"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := db.ExecContext(ctx, `DELETE FROM tasks WHERE id=$1`, task.ID); err != nil {
			t.Error(err)
		}
	})
	got, err := uc.Get(ctx, task.ID)
	if err != nil || got.ID != task.ID || got.Title != task.Title || !got.CreatedAt.Equal(task.CreatedAt) {
		t.Fatalf("read back: %+v %v", got, err)
	}
	all, err := uc.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, item := range all {
		if item.ID == task.ID {
			found = true
		}
	}
	if !found {
		t.Fatal("created task absent from list")
	}
	if _, err := uc.Get(ctx, "nonexistent-integration-task"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("missing task: %v", err)
	}
}
