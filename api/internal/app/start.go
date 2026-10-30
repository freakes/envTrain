package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"api/internal/app/server"
	"api/internal/app/tasks"
	"api/internal/app/tasks/domain"
	"api/internal/app/tasks/repository/memory"
	taskpostgres "api/internal/app/tasks/repository/postgres"
	"api/pkg/config"
	"api/pkg/postgres"

	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, cfg config.Config, log *logrus.Logger) error {
	var repo domain.Repository = memory.New()

	if cfg.Database.Driver == "postgres" {
		db, err := postgres.Open(ctx, cfg.Database.DSN)
		if err != nil {
			return err
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.WithError(err).Error("close database")
			}
		}()
		// Fail early if migrations were not applied.
		if _, err := db.ExecContext(ctx, `SELECT id,title,created_at FROM tasks LIMIT 0`); err != nil {
			return fmt.Errorf("tasks schema unavailable; apply migrations: %w", err)
		}

		repo = taskpostgres.New(db)
	}

	e := server.New(log)
	tasks.Init(e.Group("/api/v1"), repo)
	e.Server.Addr = cfg.Server.Address()
	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Server.ReadHeaderTimeout = cfg.Server.ReadTimeout
	e.Server.WriteTimeout = cfg.Server.WriteTimeout
	e.Server.IdleTimeout = cfg.Server.ReadTimeout * 6

	var listenConfig net.ListenConfig

	listener, err := listenConfig.Listen(ctx, "tcp", cfg.Server.Address())
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}

	e.Server.Handler = e

	done := make(chan error, 1)
	go func() { done <- e.Server.Serve(listener) }()

	log.WithField("address", cfg.Server.Address()).Info("starting API")

	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := e.Server.Shutdown(shutdownCtx); err != nil {
			_ = e.Server.Close()

			return fmt.Errorf("shutdown HTTP: %w", err)
		}

		err := <-done
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}

		return nil
	}
}
