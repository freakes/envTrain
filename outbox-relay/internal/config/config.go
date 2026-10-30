package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL      string
	Brokers          []string
	Topic            string
	PollInterval     time.Duration
	BatchSize        int
	OperationTimeout time.Duration
}

func value(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}
func Load() (Config, error) {
	c := Config{
		DatabaseURL: value("DATABASE_URL", "postgres://boilerplate:boilerplate@localhost:5433/events?sslmode=disable"),
		Brokers:     strings.Split(value("KAFKA_BROKERS", "localhost:19092"), ","),
		Topic:       value("KAFKA_TOPIC", "task-events"),
	}
	var err error
	c.PollInterval, err = time.ParseDuration(value("POLL_INTERVAL", "1s"))
	if err != nil || c.PollInterval <= 0 {
		return c, fmt.Errorf("POLL_INTERVAL must be a positive duration")
	}
	c.OperationTimeout, err = time.ParseDuration(value("OPERATION_TIMEOUT", "10s"))
	if err != nil || c.OperationTimeout <= 0 {
		return c, fmt.Errorf("OPERATION_TIMEOUT must be a positive duration")
	}
	c.BatchSize, err = strconv.Atoi(value("BATCH_SIZE", "10"))
	if err != nil || c.BatchSize < 1 || c.BatchSize > 100 {
		return c, fmt.Errorf("BATCH_SIZE must be between 1 and 100")
	}
	if c.DatabaseURL == "" || c.Topic == "" {
		return c, fmt.Errorf("DATABASE_URL and KAFKA_TOPIC must be nonempty")
	}
	for i, broker := range c.Brokers {
		c.Brokers[i] = strings.TrimSpace(broker)
		if c.Brokers[i] == "" {
			return c, fmt.Errorf("KAFKA_BROKERS contains an empty address")
		}
	}
	return c, nil
}
