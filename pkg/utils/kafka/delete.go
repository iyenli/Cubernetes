package kafka

import "log"

func DeleteAllTopics(host string) error {
	conn := NewKafkaClient(host)
	topics := ListTopics(host)

	for topic, _ := range topics {
		err := conn.DeleteTopics(topic)
		if err != nil {
			log.Println("[Error]: Delete all topics")
			return err
		}
	}

	return nil
}
