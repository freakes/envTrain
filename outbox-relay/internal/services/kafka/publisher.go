package kafka

import (
	"github.com/segmentio/kafka-go"
)

type Publisher struct{ writer *kafka.Writer }
