package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

const TaskCreated = "task.created.v1"

type Event struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	AggregateID string          `json:"aggregate_id"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Payload     json.RawMessage `json:"payload"`
}
type Task struct {
	ID    string `json:"task_id"`
	Title string `json:"title"`
}

func (e Event) Task() (Task, error) {
	if strings.TrimSpace(e.ID) == "" || e.Type != TaskCreated || strings.TrimSpace(e.AggregateID) == "" || e.OccurredAt.IsZero() {
		return Task{}, errors.New("invalid or unsupported event envelope")
	}
	var task Task
	if err := json.Unmarshal(e.Payload, &task); err != nil {
		return task, errors.New("invalid event payload")
	}
	if task.ID != e.AggregateID || strings.TrimSpace(task.Title) == "" || utf8.RuneCountInString(task.Title) > 200 {
		return task, errors.New("invalid task payload")
	}
	return task, nil
}
