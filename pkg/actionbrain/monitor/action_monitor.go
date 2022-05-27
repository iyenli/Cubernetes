package monitor

import (
	"Cubernetes/pkg/actionbrain/monitor/message"
	"Cubernetes/pkg/actionbrain/monitor/options"
	kafkautil "Cubernetes/pkg/utils/kafka"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nakabonne/tstorage"
	"github.com/segmentio/kafka-go"
)

type ActionMonitor interface {
	Run()
	QueryRecentEvoke(action string, period time.Duration) (int, error)
	WatchActionEvoke() <-chan string
	Close()
}

func NewActionMonitor() (ActionMonitor, error) {
	reader := kafkautil.NewReaderByConsumerGroup(
		options.KafkaHost, options.MonitorTopic, options.MonitorConsumerGroup)

	storage, err := tstorage.NewStorage(tstorage.WithPartitionDuration(30 * time.Second))
	if err != nil {
		log.Printf("fail to create in-memory storage: %v\n", err)
		return nil, err
	}

	evokeChan := make(chan string, options.ActionChanBufferLen)

	return &kafkaActionMonitor{
		reader:    reader,
		storage:   storage,
		evokeChan: evokeChan,
	}, nil
}

type kafkaActionMonitor struct {
	reader    *kafka.Reader
	storage   tstorage.Storage
	evokeChan chan string
}

func (kam *kafkaActionMonitor) Run() {
	ctx := context.Background()
	for {
		msgByte, err := kam.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("fail to read message: %v", err)
			break
		}

		msg := message.MonitorMessage{}
		if err = json.Unmarshal(msgByte.Value, &msg); err != nil {
			log.Printf("fail to parse MonitorMessage: %s", string(msgByte.Value))
			continue
		}

		action := msg.Action
		time := msg.InvokeTimeUnix
		if err = kam.insertActionEvoke(action, time); err != nil {
			log.Printf("fail to storage action evoke: %v", err)
		}
		kam.evokeChan <- action
	}
}

func (kam *kafkaActionMonitor) insertActionEvoke(action string, time int64) error {
	return kam.storage.InsertRows([]tstorage.Row{
		{
			Metric:    "action",
			Labels:    []tstorage.Label{{Name: "action_name", Value: action}},
			DataPoint: tstorage.DataPoint{Timestamp: time, Value: 0.0},
		},
	})
}

func (kam *kafkaActionMonitor) QueryRecentEvoke(action string, period time.Duration) (int, error) {
	now := time.Now().Unix()
	points, err := kam.storage.Select(
		"action",
		[]tstorage.Label{{Name: "action_name", Value: action}},
		now-int64(period.Seconds()), now,
	)
	if err != nil {
		log.Printf("fail to query evoke of action %s: %v", action, err)
		return 0, err
	}

	return len(points), nil
}

func (kam *kafkaActionMonitor) WatchActionEvoke() <-chan string {
	return kam.evokeChan
}

func (kam *kafkaActionMonitor) Close() {
	kam.reader.Close()
	kam.storage.Close()
	close(kam.evokeChan)
}
