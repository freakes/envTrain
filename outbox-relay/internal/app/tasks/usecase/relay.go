package usecase

import (
	"relay/internal/app/tasks/domain"
)

type Relay struct {
	repo      domain.Repository
	publisher domain.Publisher
	batchSize int
}

func New(repo domain.Repository, publisher domain.Publisher, batchSize int) *Relay {
	return &Relay{repo: repo, publisher: publisher, batchSize: batchSize}
}
