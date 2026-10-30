package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func New(level string, json bool) (*logrus.Logger, error) {
	parsed, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	log := logrus.New()
	log.SetLevel(parsed)

	if json {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	return log, nil
}
