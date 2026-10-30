package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Event is a transport-neutral envelope. The relay does not interpret business payloads.
type Event struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	AggregateID string          `json:"aggregate_id"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Payload     json.RawMessage `json:"payload"`
}

func (e Event) Validate() error {
	if strings.TrimSpace(e.ID) == "" || strings.TrimSpace(e.Type) == "" || strings.TrimSpace(e.AggregateID) == "" || e.OccurredAt.IsZero() || !json.Valid(e.Payload) {
		return errors.New("invalid event envelope")
	}
	return nil
}
