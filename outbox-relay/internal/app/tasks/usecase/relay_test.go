package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"relay/internal/app/tasks/domain"
	"relay/internal/app/tasks/usecase"
)

type repository struct {
	batch     *batch
	committed bool
}

func (r *repository) WithBatch(ctx context.Context, limit int, fn func(domain.Batch) error) error {
	err := fn(r.batch)
	r.committed = err == nil
	return err
}

type batch struct {
	events  []domain.Event
	calls   *[]string
	markErr error
}

func (b *batch) Events() []domain.Event { return b.events }
func (b *batch) MarkPublished(ctx context.Context, id string) error {
	*b.calls = append(*b.calls, "mark:"+id)
	return b.markErr
}

type publisher struct {
	calls *[]string
	err   error
}

func (p publisher) Publish(ctx context.Context, e domain.Event) error {
	*p.calls = append(*p.calls, "publish:"+e.ID)
	return p.err
}
func TestPublishBeforeMarkAndRollback(t *testing.T) {
	failure := errors.New("unavailable")
	for _, tt := range []struct {
		name                string
		publishErr, markErr error
		want                []string
		committed           bool
	}{
		{"success", nil, nil, []string{"publish:a", "mark:a", "publish:b", "mark:b"}, true},
		{"publish failed", failure, nil, []string{"publish:a"}, false},
		{"mark failed", nil, failure, []string{"publish:a", "mark:a"}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			calls := []string{}
			events := []domain.Event{}
			for _, id := range []string{"a", "b"} {
				events = append(events, domain.Event{ID: id, Type: "any.event.v1", AggregateID: "key", OccurredAt: time.Now(), Payload: json.RawMessage(`{}`)})
			}
			repo := &repository{batch: &batch{events: events, calls: &calls, markErr: tt.markErr}}
			err := usecase.New(repo, publisher{calls: &calls, err: tt.publishErr}, 10).RunOnce(context.Background())
			if tt.committed && err != nil {
				t.Fatal(err)
			}
			if !tt.committed && !errors.Is(err, failure) {
				t.Fatalf("lost cause: %v", err)
			}
			if repo.committed != tt.committed || !reflect.DeepEqual(calls, tt.want) {
				t.Fatalf("committed=%v calls=%v", repo.committed, calls)
			}
		})
	}
}
func TestInvalidEnvelopeIsNotPublished(t *testing.T) {
	calls := []string{}
	repo := &repository{batch: &batch{events: []domain.Event{{ID: "broken"}}, calls: &calls}}
	if err := usecase.New(repo, publisher{calls: &calls}, 1).RunOnce(context.Background()); err == nil {
		t.Fatal("want validation failure")
	}
	if len(calls) != 0 || repo.committed {
		t.Fatal("invalid event published")
	}
}
