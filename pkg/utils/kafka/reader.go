package kafka

import (
	"Cubernetes/pkg/utils/kafka/options"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
)

// NewReader
/**
  Usage:
	one phase:
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			// log here:)
		}
	}

  Or two phase
	ctx := context.Background()
	for {
    	m, err := r.FetchMessage(ctx)
		if err != nil {
     	   // log here:)
   	 	}
		// consume it
    	if err := r.CommitMessages(ctx, m); err != nil {
        	// log here:)
    	}
	}
*/
func NewReader(host string, topic string, partition int, offset int64) *kafka.Reader {
	address := net.JoinHostPort(host, options.KafkaPort)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     topic,
		Partition: partition,
		MinBytes:  options.ReaderMinByte,
		MaxBytes:  options.ReaderMaxByte,

		Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	})

	err := r.SetOffset(offset)
	if err != nil {
		log.Printf("[Error]: Create new reader with topic %v, partition %v failed", topic, partition)
		return nil
	}

	return r
}

// NewReaderByConsumerGroup More recommended, Kafka helps you manage offset
func NewReaderByConsumerGroup(host, topic, consumerGroupID string) *kafka.Reader {
	address := net.JoinHostPort(host, options.KafkaPort)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
		GroupID: consumerGroupID,

		MinBytes: options.ReaderMinByte,
		MaxBytes: options.ReaderMaxByte,

		Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	})

	return r
}
