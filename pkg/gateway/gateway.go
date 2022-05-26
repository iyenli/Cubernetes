package gateway

import (
	"Cubernetes/pkg/gateway/httpserver"
	"Cubernetes/pkg/gateway/informer"
	"Cubernetes/pkg/gateway/options"
	"Cubernetes/pkg/gateway/types"
	kafka2 "Cubernetes/pkg/utils/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
)

type RuntimeGateway struct {
	router          *gin.Engine
	ingressInformer informer.IngressInformer

	channelMap map[string]chan types.MQMessage
	mapMutex   sync.Mutex

	returnTopic string

	writer *kafka.Writer
	reader *kafka.Reader
}

func NewRuntimeGateway() *RuntimeGateway {
	returnTopic := options.ListenTopicPrefix + uuid.NewString()
	// Create relevant topic
	err := kafka2.CreateTopic("127.0.0.1", returnTopic)
	if err != nil {
		log.Println("[Error]: create topic failed")
		return nil
	}

	consumerGroupID := "GatewayConsumerGroup-" + returnTopic
	return &RuntimeGateway{
		router:          httpserver.GetGatewayRouter(),
		ingressInformer: informer.NewIngressInformer(),

		channelMap: make(map[string]chan types.MQMessage),
		mapMutex:   sync.Mutex{},

		returnTopic: returnTopic,

		writer: kafka2.NewWriter("127.0.0.1"),
		reader: kafka2.NewReaderByConsumerGroup("127.0.0.1", returnTopic, consumerGroupID),
	}
}

func (rg *RuntimeGateway) Run() {
	wg := sync.WaitGroup{}
	wg.Add(4)

	// Gateway
	go func() {
		defer wg.Done()
		httpserver.Run(rg.router)
	}()

	// RETURN TOPIC Consumer
	go func() {
		defer wg.Done()
		rg.ListenReturnTopic()
	}()

	// Ingress informer
	go func() {
		defer wg.Done()
		rg.ingressInformer.ListAndWatchIngressWithRetry()
	}()

	// Handler with ingress CRUD
	go func() {
		defer wg.Done()
		rg.HandleIngress()
	}()

	log.Println("[INFO]: Gateway running...")
	wg.Wait()
	log.Println("[Fatal]: Gateway exited incorrectly")
}
