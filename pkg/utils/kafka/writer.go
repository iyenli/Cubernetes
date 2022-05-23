package kafka

import (
	"Cubernetes/pkg/utils/kafka/options"
	"github.com/segmentio/kafka-go"
	"net"
)

// NewWriterWithTopic it'll create the topic if topic hasn't created
/**
Usage:
err := w.WriteMessages(context.Background(),
	kafka.Message{
		Key:   []byte("Key-A"),
		Value: []byte("Hello World!"),
	}, (...More msg is okay)
)
*/
func NewWriterWithTopic(host string, topic string) *kafka.Writer {
	address := net.JoinHostPort(host, options.KafkaPort)
	w := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},

		Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	}

	return w
}

// NewWriter
/**
Usage:
err := w.WriteMessages(context.Background(),
    // NOTE: Each Message has Topic defined, otherwise an error is returned.
	kafka.Message{
        Topic: "topic-A",
		Key:   []byte("Key-A"),
		Value: []byte("Hello World!"),
	},
)
*/
func NewWriter(host string) *kafka.Writer {
	address := net.JoinHostPort(host, options.KafkaPort)

	w := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Balancer: &kafka.LeastBytes{},

		Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	}

	return w
}
