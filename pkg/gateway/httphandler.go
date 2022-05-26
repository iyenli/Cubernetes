package gateway

import (
	"Cubernetes/pkg/gateway/types"
	"Cubernetes/pkg/object"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"strconv"
)

func (rg *RuntimeGateway) GetHandlerByIngress(ingress *object.Ingress) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		log.Printf("[INFO]: Function called, ingress path: %v, input topic %v",
			ingress.Spec.TriggerPath, ingress.Spec.InvokeAction)

		msg := types.MQMessage{
			RequestUID:  uuid.NewString(),
			TriggerPath: ingress.Spec.TriggerPath,
			ReturnTopic: rg.returnTopic,
			ContentType: "", // wait for filling
			Params:      make(map[string]string),
			Payload:     "",
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

		if ingress.Spec.HTTPType != http.MethodGet && ctx.ContentType() == gin.MIMEJSON {
			body := make(map[string]interface{})

			err := ctx.BindJSON(&body)
			if err != nil {
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
			msg.Payload = string(byteBody)
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[INFO]: Marshal message failed. Trigger path %v, id %v, body %v",
				msg.TriggerPath, msg.RequestUID, msg.Payload)
			return
		}

		channel := make(chan types.MQMessage, 1)
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
			ctx.String(http.StatusInternalServerError, "Kafka Error")
			return
		}

		// Just wait in go routine...
		log.Printf("[INFO]: Request %v has been put into MQ, waiting for resp...\n", msg.RequestUID)
		resp := <-channel
		log.Printf("[INFO]: Req has got the resp! ID is %v", msg.RequestUID)
		code, err := strconv.Atoi(resp.StatusCode)
		if err != nil {
			log.Printf("[Error]: return code is invalid")
			ctx.String(http.StatusInternalServerError, "HttpCode error")
			return
		}

		log.Printf("[INFO]: Resp received, body is %v, return type is %v",
			resp.Payload, resp.ContentType)
		ctx.Header("Content-Type", resp.ContentType)
		ctx.String(code, resp.Payload)
	}
}
