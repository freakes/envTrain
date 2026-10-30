package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"api/internal/app"
	"api/pkg/config"
	"api/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.Getenv("CONFIG_PATH"))
	if err != nil {
		return err
	}

	log, err := logger.New(cfg.Log.Level, cfg.Log.JSON)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return app.Run(ctx, cfg, log)
}
