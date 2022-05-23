package kafka

import (
	"Cubernetes/pkg/utils/kafka/options"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
)

func CreateTopic(host string, topic string) error {
	conn := NewKafkaClient(host)
	if conn == nil {
		log.Println("[Fatal]: Create conn failed")
		return errors.New("get Kafka client failed")
	}

	defer func(controllerConn *kafka.Conn) {
		err := controllerConn.Close()
		if err != nil {
			log.Println("[Error]: close connection with master failed")
		}
	}(conn)

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     options.MaxFunctionReplica,
			ReplicationFactor: 1, // unset replication
		},
	}

	err := conn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Println("[Error]: create topic failed")
		return err
	}

	return nil
}

func ListTopics(host string) (mp map[string]struct{}) {
	conn := NewKafkaClient(host)
	if conn == nil {
		log.Println("[Fatal]: Create conn failed")
		return
	}

	defer func(controllerConn *kafka.Conn) {
		err := controllerConn.Close()
		if err != nil {
			log.Println("[Error]: close connection with master failed")
		}
	}(conn)

	partitions, err := conn.ReadPartitions()
	if err != nil {
		log.Println("[Fatal]: Get Partition failed")
		return
	}

	for _, p := range partitions {
		mp[p.Topic] = struct{}{}
	}

	return
}
