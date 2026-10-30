package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"api/pkg/config"
)

func TestDefaultsAndEnv(t *testing.T) {
	t.Setenv("APP_SERVER_PORT", "9090")
	t.Setenv("APP_SERVER_READ_TIMEOUT", "3s")
	c, err := config.Load("")
	if err != nil {
		t.Fatal(err)
	}
	if c.Server.Port != 9090 || c.Server.ReadTimeout != 3*time.Second || c.Database.Driver != "memory" {
		t.Fatalf("bad config: %+v", c)
	}
}
func TestFileAndEnvPrecedence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"server":{"port":8081},"log":{"level":"debug"}}`), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_SERVER_PORT", "8082")
	c, err := config.Load(path)
	if err != nil || c.Server.Port != 8082 || c.Log.Level != "debug" {
		t.Fatalf("bad config: %+v %v", c, err)
	}
}
func TestInvalidConfiguration(t *testing.T) {
	for _, tt := range []struct{ key, value string }{
		{"APP_SERVER_PORT", "0"}, {"APP_SERVER_PORT", "no"}, {"APP_SERVER_READ_TIMEOUT", "0s"},
		{"APP_DATABASE_DRIVER", "sqlite"}, {"APP_DATABASE_DRIVER", "postgres"},
	} {
		t.Run(tt.key+tt.value, func(t *testing.T) {
			t.Setenv(tt.key, tt.value)
			if _, err := config.Load(""); err == nil {
				t.Fatal("expected error")
			}
		})
	}
	if _, err := config.Load("/missing/config.json"); err == nil {
		t.Fatal("explicit missing file must fail")
	}
}
