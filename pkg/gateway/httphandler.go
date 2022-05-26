package gateway

import (
	"Cubernetes/pkg/object"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
)

func (rg *RuntimeGateway) GetHandlerByIngress(ingress *object.Ingress) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		log.Printf("[INFO]: Function called, ingress path: %v, input topic %v",
			ingress.Spec.TriggerPath, ingress.Spec.InvokeAction)

		msg := object.MQMessage{
			RequestUID:  uuid.NewString(),
			TriggerPath: ingress.Spec.TriggerPath,
			ReturnTopic: rg.returnTopic,
			ReturnType:  "", // wait for filling
			Params:      make(map[string]string),
			Body:        "",
		}

		for key, val := range ctx.Request.URL.Query() {
			if len(val) == 0 {
				continue // ignore the param
			} else if len(val) > 1 {
				log.Printf("[Error]: One param has more than two values is not allowed")
				ctx.String(http.StatusBadRequest, "One param has more than two values is not allowed")
				return
			}
			log.Printf("[INFO]: Get param key is %v, value is %v", key, val[0])
			msg.Params[key] = val[0]
		}
		if ctx.ContentType() != http.MethodGet {
			body := make(map[string]interface{})

			switch ctx.ContentType() {
			case gin.MIMEJSON:
				err := ctx.BindJSON(&body)
				if err != nil {
					log.Printf("[Error]: Just accept json as Content-Type")
					ctx.String(http.StatusBadRequest, "JSON is required")
					return
				}

			default:
				log.Printf("[Error]: Just accept json as Content-Type")
				ctx.String(http.StatusBadRequest, "JSON is required")
				return
			}

			byteBody, err := json.Marshal(body)
			if err != nil {
				log.Printf("[Error]: Just accept json as Content-Type")
				ctx.String(http.StatusBadRequest, "JSON is required")
				return
			}

			log.Printf("[INFO]: Http JSON body is %v", string(byteBody))
			msg.Body = string(byteBody)
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

		log.Printf("[INFO]: Resp received, body is %v, return type is %v",
			resp.Body, resp.ReturnType)
		ctx.Header("Content-Type", resp.ReturnType)
		ctx.String(http.StatusOK, resp.Body)
	}
}
