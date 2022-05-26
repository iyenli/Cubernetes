package gateway

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	kafka2 "Cubernetes/pkg/utils/kafka"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

func (rg *RuntimeGateway) NewConsumer() *kafka.Reader {
	consumerGroupID := "GatewayConsumerGroup-" + rg.returnTopic
	return kafka2.NewReaderByConsumerGroup(cubeconfig.APIServerIp, rg.returnTopic, consumerGroupID)
}

func (rg *RuntimeGateway) ListenReturnTopic() {
	ctx := context.Background()
	for {
		msgByte, err := rg.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("[Error]: read message failed")
			continue
		}

		msg := object.MQMessage{}
		err = json.Unmarshal(msgByte.Value, &msg)
		if err != nil {
			log.Printf("[Error]: parse return message failed, msg: %v", string(msgByte.Value))
			return
		}

		rg.mapMutex.Lock()
		if channel, ok := rg.channelMap[msg.RequestUID]; ok {
			log.Printf("[INFO]: Get resp from MQ, UID is %v", msg.RequestUID)
			channel <- msg
		} else {
			log.Printf("[Error]: return a resp and not found in channel map")
		}
		rg.mapMutex.Unlock()
	}
}
