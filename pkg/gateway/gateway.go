package gateway

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/gateway/httpserver"
	"Cubernetes/pkg/gateway/informer"
	"Cubernetes/pkg/gateway/informer/types"
	"Cubernetes/pkg/gateway/options"
	"Cubernetes/pkg/object"
	kafka2 "Cubernetes/pkg/utils/kafka"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"sync"
)

type RuntimeGateway struct {
	router          *gin.Engine
	ingressInformer informer.IngressInformer

	channelMap map[string]chan object.MQMessage
	mapMutex   sync.Mutex

	returnTopic string

	writer *kafka.Writer
	reader *kafka.Reader
}

func NewRuntimeGateway() *RuntimeGateway {
	returnTopic := options.TopicPrefix + uuid.NewString()
	return &RuntimeGateway{
		router:          httpserver.GetGatewayRouter(),
		ingressInformer: informer.NewIngressInformer(),

		channelMap: make(map[string]chan object.MQMessage),
		mapMutex:   sync.Mutex{},

		returnTopic: returnTopic,

		writer: kafka2.NewWriter(cubeconfig.APIServerIp),
	}
}

func (rg *RuntimeGateway) Run() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		httpserver.Run(rg.router)
	}()

	go func() {
		defer wg.Done()
		rg.ingressInformer.ListAndWatchIngressWithRetry()
	}()

	go func() {
		defer wg.Done()
		rg.HandleIngress()
	}()

	log.Println("[INFO]: Gateway running...")
	wg.Wait()
	log.Println("[Fatal]: Gateway exited incorrectly")
}

func (rg *RuntimeGateway) HandleIngress() {
	informEvent := rg.ingressInformer.WatchIngressEvent()

	for ingressEvent := range informEvent {
		log.Printf("[INFO]: Main loop working, types is %v, ingress id is %v",
			ingressEvent.Type, ingressEvent.Ingress.UID)

		switch ingressEvent.Type {
		case types.IngressCreate:
			log.Printf("[INFO]: Create Ingress, UID is %v", ingressEvent.Ingress.UID)
			// Add to router
			rg.AddIngress(&ingressEvent.Ingress)
		case types.IngressUpdate:
			log.Printf("[INFO]: Update Ingress, UID is %v", ingressEvent.Ingress.UID)
		case types.IngressRemove:
			log.Printf("[INFO]: Delete Ingress, UID is %v", ingressEvent.Ingress.UID)
		}
	}
}

func (rg *RuntimeGateway) AddIngress(ingress *object.Ingress) {
	switch ingress.Spec.HTTPType {
	case http.MethodPut:
		rg.router.PUT(ingress.Spec.TriggerPath)
	case http.MethodGet:
		rg.router.GET(ingress.Spec.TriggerPath)
	case http.MethodDelete:
		rg.router.DELETE(ingress.Spec.TriggerPath)
	case http.MethodPost:
		rg.router.POST(ingress.Spec.TriggerPath)
	default:
		log.Printf("[Warn]: unsupported type of http, discard it")
		return

	}
}

func (rg *RuntimeGateway) GetHandlerByIngress(ingress *object.Ingress, body bool, uri bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		log.Printf("[INFO]: Function called, ingress path: %v, input topic %v",
			ingress.Spec.TriggerPath, ingress.Spec.InvokeAction)

		msg := object.MQMessage{
			RequestUID:  uuid.NewString(),
			ReturnTopic: rg.returnTopic,
			TriggerPath: ingress.Spec.TriggerPath,
			Body:        "",
		}
		if body {
			switch ctx.ContentType() {
			case gin.MIMEJSON:

			case gin.MIMEPOSTForm:

			case gin.MIMEPlain:

			default:

			}
		}
		if uri {

		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[INFO]: Marshal message failed. Trigger path %v, id %v, body %v",
				msg.TriggerPath, msg.RequestUID, msg.Body)
			return
		}

		channel := make(chan object.MQMessage, 1)
		rg.mapMutex.Lock()
		rg.channelMap[msg.RequestUID] = channel
		rg.mapMutex.Unlock()

		err = rg.writer.WriteMessages(context.Background(),
			kafka.Message{
				Topic: ingress.Spec.InvokeAction,
				Key:   []byte(msg.RequestUID),
				Value: msgBytes,
			})
		if err != nil {
			log.Printf("[Error]: Write into MQ %v failed", ingress.Spec.InvokeAction)
			return
		}

		// Just wait in go routine...
		log.Printf("[INFO]: Request %v has been put into MQ, waiting for resp...\n", msg.RequestUID)
		resp := <-channel
		ctx.String(http.StatusOK, resp.Body)
	}
}
