package kafka

import (
	"Cubernetes/pkg/utils/kafka/options"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
)

func CreateTopic(conn *kafka.Conn, topic string) error {
	controller, err := conn.Controller()
	if err != nil {
		log.Println("[Error]: get controller failed")
		return err
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Println("[Error]: get controller failed")
		return err
	}

	defer func(controllerConn *kafka.Conn) {
		err := controllerConn.Close()
		if err != nil {
			log.Println("[Error]: close controller failed")
		}
	}(controllerConn)

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     options.MaxFunctionReplica,
			ReplicationFactor: 1, // unset replication
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Println("[Error]: create topic failed")
		return err
	}
	return nil
}
