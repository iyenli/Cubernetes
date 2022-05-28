package testing

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/utils/kafka"
	"context"
	kafka2 "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

// Give it a try:)
// go test ./pkg/utils/kafka/testing/kafkaUtils_test.go -v -test.run TestNewKafkaClient

const (
	HOST = "127.0.0.1"
)

func TestNewKafkaClient(t *testing.T) {
	topic := "test-topic"
	err := kafka.CreateTopic(HOST, topic)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(6)

	writer := kafka.NewWriter(HOST)
	assert.NotNil(t, writer)

	readers := make([]*kafka2.Reader, 5)
	for i := 0; i < 5; i++ {
		readers[i] = kafka.NewReaderByConsumerGroup(HOST, topic, topic)
		idx := i
		go func(idx int) {
			defer wg.Done()
			for {
				message, err := readers[idx].ReadMessage(context.Background())
				assert.NoError(t, err)
				log.Printf("[INFO]: msg key %v, msg val %v, index %v",
					string(message.Key), string(message.Value), idx)
			}
		}(idx)
	}

	go func() {
		defer wg.Done()
		for {
			err := writer.WriteMessages(context.Background(), kafka2.Message{
				Topic: topic,
				Key:   []byte(time.Now().String()),
				Value: []byte(time.Now().String()),
			})
			assert.NoError(t, err)

			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
}

// go test pkg/utils/kafka/testing/kafkaUtils_test.go -test.run TestClean
func TestClean(t *testing.T) {
	err := kafka.DeleteAllTopics(cubeconfig.APIServerIp)
	assert.NoError(t, err)
}
