package app_test

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"api/internal/app"
	"api/pkg/config"
	"github.com/sirupsen/logrus"
)

func TestStartupFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	cfg, err := config.Load("")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = listener.Addr().(*net.TCPAddr).Port
	log := logrus.New()
	log.SetOutput(io.Discard)
	if err := app.Run(context.Background(), cfg, log); err == nil {
		t.Fatal("occupied port must fail startup")
	}
}

func TestCanceledStartupDoesNotHang(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 0 // OS-assigned port for this lifecycle test; public config rejects zero.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	log := logrus.New()
	log.SetOutput(io.Discard)
	done := make(chan error, 1)
	go func() { done <- app.Run(ctx, cfg, log) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("shutdown hung")
	}
}
