package config_test

import (
	"testing"

	"relay/internal/config"
)

func TestConfigValidation(t *testing.T) {
	for _, tt := range []struct{ key, value string }{
		{"BATCH_SIZE", "0"}, {"BATCH_SIZE", "101"}, {"POLL_INTERVAL", "0s"},
		{"OPERATION_TIMEOUT", "bad"}, {"KAFKA_BROKERS", "host:9092,"}, {"KAFKA_TOPIC", ""}, {"DATABASE_URL", ""},
	} {
		t.Run(tt.key+tt.value, func(t *testing.T) {
			t.Setenv(tt.key, tt.value)
			if _, err := config.Load(); err == nil {
				t.Fatal("want config error")
			}
		})
	}
}
func TestDefaultsAndOverride(t *testing.T) {
	t.Setenv("BATCH_SIZE", "3")
	c, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.BatchSize != 3 || len(c.Brokers) == 0 {
		t.Fatalf("config: %+v", c)
	}
}
