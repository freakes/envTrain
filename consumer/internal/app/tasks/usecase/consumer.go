package usecase

import (
	"context"
	"fmt"
	"outbox-consumer/internal/app/tasks/domain"
)

type Consumer struct{ repo domain.Repository }

func New(repo domain.Repository) *Consumer { return &Consumer{repo: repo} }
func (c *Consumer) Handle(ctx context.Context, event domain.Event) error {

}
