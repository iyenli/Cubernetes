package monitor

import (
	"Cubernetes/pkg/actionbrain/monitor/message"
	"Cubernetes/pkg/actionbrain/monitor/options"
	kafkautil "Cubernetes/pkg/utils/kafka"
	"context"
	"encoding/json"
	"errors"
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

func NewActionMonitor(kafkaHost string) (ActionMonitor, error) {

	if err := kafkautil.CreateTopic(kafkaHost, options.MonitorTopic); err != nil {
		log.Printf("fail to create monitor topic: %v\n", err)
		return nil, err
	}

	reader := kafkautil.NewReaderByConsumerGroup(
		kafkaHost, options.MonitorTopic, options.MonitorConsumerGroup)

	storage, err := tstorage.NewStorage(
		tstorage.WithPartitionDuration(30000 * time.Second),
	)
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
		log.Printf("received invoke to action %s\n", action)
		if err = kam.insertActionEvoke(action, time); err != nil {
			log.Printf("fail to storage action evoke: %v", err)
		}
		kam.evokeChan <- action
	}
}

func (kam *kafkaActionMonitor) insertActionEvoke(action string, t int64) error {
	return kam.storage.InsertRows([]tstorage.Row{
		{
			Metric:    "action",
			Labels:    []tstorage.Label{{Name: "action_name", Value: action}},
			DataPoint: tstorage.DataPoint{Value: 0.0, Timestamp: time.Now().UnixMilli()},
		},
	})
}

func (kam *kafkaActionMonitor) QueryRecentEvoke(action string, period time.Duration) (int, error) {
	now := time.Now().UnixMilli()
	points, err := kam.storage.Select(
		"action",
		[]tstorage.Label{{Name: "action_name", Value: action}},
		now-period.Milliseconds(), now,
	)
	if err != nil {
		if errors.Is(err, tstorage.ErrNoDataPoints) {
			return 0, nil
		}
		log.Printf("fail to query evoke of action %s: %v", action, err)
		return 0, err
	}
	log.Printf("total %d evoke(s) for action %s found\n", len(points), action)

	return len(points), nil
}

func (kam *kafkaActionMonitor) WatchActionEvoke() <-chan string {
	return kam.evokeChan
}

func (kam *kafkaActionMonitor) Close() {
	err := kam.reader.Close()
	if err != nil {
		return
	}
	err = kam.storage.Close()
	if err != nil {
		return
	}
	close(kam.evokeChan)
}
