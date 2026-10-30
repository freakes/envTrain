package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL      string
	Brokers          []string
	Topic            string
	Group            string
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
		Topic:       value("KAFKA_TOPIC", "task-events"), Group: value("KAFKA_GROUP", "outbox-consumer"),
	}
	var err error
	c.OperationTimeout, err = time.ParseDuration(value("OPERATION_TIMEOUT", "10s"))
	if err != nil || c.OperationTimeout <= 0 {
		return c, fmt.Errorf("OPERATION_TIMEOUT must be a positive duration")
	}
	if c.DatabaseURL == "" || c.Topic == "" || c.Group == "" {
		return c, fmt.Errorf("DATABASE_URL, KAFKA_TOPIC and KAFKA_GROUP must be nonempty")
	}
	for i, broker := range c.Brokers {
		c.Brokers[i] = strings.TrimSpace(broker)
		if c.Brokers[i] == "" {
			return c, fmt.Errorf("KAFKA_BROKERS contains an empty address")
		}
	}
	return c, nil
}
