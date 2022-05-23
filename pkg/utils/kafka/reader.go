package kafka

import (
	"Cubernetes/pkg/utils/kafka/options"
	"github.com/segmentio/kafka-go"
	"net"
)

// NewReader
/**
  Usage:
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			// log here:)
		}
	}
*/
func NewReader(host string, topic string, partition int) *kafka.Reader {
	address := net.JoinHostPort(host, options.KafkaPort)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	err := r.SetOffset(0)
	if err != nil {
		return nil
	}

	return r
}
