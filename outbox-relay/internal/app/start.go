package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	repository "relay/internal/app/tasks/repository/postgres"
	"relay/internal/app/tasks/usecase"
	"relay/internal/config"
	"relay/internal/services/kafka"
	"relay/internal/storage"
)

func Run(ctx context.Context, cfg config.Config, log *slog.Logger) error {

}
