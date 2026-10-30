package postgres_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"relay/internal/app/tasks/domain"
	repository "relay/internal/app/tasks/repository/postgres"
	"relay/internal/storage"
)

func TestLockSkipRollbackAndCommit(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL for PostgreSQL integration")
	}
	ctx := context.Background()
	db, err := storage.Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	id := rand.Text()
	event := domain.Event{ID: id, Type: "test.event.v1", AggregateID: id, OccurredAt: time.Now(), Payload: json.RawMessage(`{}`)}
	raw, _ := json.Marshal(event)
	if _, err := db.ExecContext(ctx, `INSERT INTO outbox_events (id,event) VALUES ($1,$2)`, id, string(raw)); err != nil {
		t.Fatal(err)
	}
	defer db.ExecContext(ctx, `DELETE FROM outbox_events WHERE id=$1`, id)
	// Lock the test row; a second relay must skip it without waiting.
	lock, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Rollback()
	if _, err := lock.ExecContext(ctx, `SELECT id FROM outbox_events WHERE id=$1 FOR UPDATE`, id); err != nil {
		t.Fatal(err)
	}
	repo := repository.New(db)
	short, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := repo.WithBatch(short, 100, func(b domain.Batch) error {
		for _, e := range b.Events() {
			if e.ID == id {
				t.Error("locked row claimed by another relay")
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := lock.Rollback(); err != nil {
		t.Fatal(err)
	}
	failure := errors.New("broker failed after ack")
	found := false
	if err := repo.WithBatch(ctx, 100, func(b domain.Batch) error {
		for _, e := range b.Events() {
			if e.ID == id {
				found = true
				if err := b.MarkPublished(ctx, id); err != nil {
					return err
				}
			}
		}
		return failure
	}); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("test row not selected; use an isolated test database")
	}
	var pending bool
	if err := db.QueryRowContext(ctx, `SELECT published_at IS NULL FROM outbox_events WHERE id=$1`, id).Scan(&pending); err != nil || !pending {
		t.Fatalf("mark survived rollback: %v %v", pending, err)
	}
	if err := repo.WithBatch(ctx, 100, func(b domain.Batch) error {
		for _, e := range b.Events() {
			if e.ID == id {
				return b.MarkPublished(ctx, id)
			}
		}
		return errors.New("row not selected")
	}); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, `SELECT published_at IS NULL FROM outbox_events WHERE id=$1`, id).Scan(&pending); err != nil || pending {
		t.Fatalf("commit lost mark: %v %v", pending, err)
	}
}
