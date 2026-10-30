package memory_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"api/internal/app/tasks/domain"
	"api/internal/app/tasks/repository/memory"
)

func TestConcurrentAccessAndOrdering(t *testing.T) {
	repo := memory.New()
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task := domain.Task{ID: fmt.Sprintf("%03d", i), Title: "example", CreatedAt: time.Unix(0, 0)}
			if err := repo.Create(ctx, task); err != nil {
				t.Error(err)
			}
			if _, err := repo.Get(ctx, task.ID); err != nil {
				t.Error(err)
			}
			if _, err := repo.List(ctx); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	all, err := repo.List(ctx)
	if err != nil || len(all) != 40 {
		t.Fatalf("list: %d %v", len(all), err)
	}
	for i, task := range all {
		if task.ID != fmt.Sprintf("%03d", i) {
			t.Fatal("unstable order")
		}
	}
	all[0].Title = "changed outside repository"
	stored, err := repo.Get(ctx, all[0].ID)
	if err != nil || stored.Title != "example" {
		t.Fatalf("returned mutable storage: %+v %v", stored, err)
	}
}
