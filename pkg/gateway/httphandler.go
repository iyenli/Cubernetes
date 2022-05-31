package gateway

import (
	"Cubernetes/pkg/actionbrain/monitor/message"
	"Cubernetes/pkg/actionbrain/monitor/options"
	"Cubernetes/pkg/gateway/types"
	"Cubernetes/pkg/gateway/utils"
	"Cubernetes/pkg/object"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"time"
)

func (rg *RuntimeGateway) GetHandlerByIngress(ingress *object.Ingress) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		log.Printf("[INFO]: Function called, ingress path: %v, input topic %v",
			ingress.Spec.TriggerPath, ingress.Spec.InvokeAction)

		// Send information to action monitor
		go func(action string) {
			err := rg.SendMonitorInfo(action)
			if err != nil {
				log.Printf("[Error]: send info to monitor failed\n")
				return
			}
		}(ingress.Spec.InvokeAction)

		msg := types.MQMessageRequest{
			RequestUID:  uuid.NewString(),
			TriggerPath: ingress.Spec.TriggerPath,
			ReturnTopic: rg.returnTopic,
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

		channel := make(chan types.MQMessageResponse, 1)
		rg.mapMutex.Lock()
		rg.channelMap[msg.RequestUID] = channel
		rg.mapMutex.Unlock()

		err = rg.writer.WriteMessages(context.Background(),
			kafka.Message{
				Topic: utils.GetActionTopic(ingress.Spec.InvokeAction),
				Key:   []byte("\"" + msg.RequestUID + "\""),
				Value: msgBytes,
			})
		if err != nil {
			log.Printf("[Error]: Write into MQ %v failed", utils.GetActionTopic(ingress.Spec.InvokeAction))
			ctx.String(http.StatusInternalServerError, "Kafka Error")
			return
		}

		// Just wait in go routine...
		log.Printf("[INFO]: Request %v has been put into MQ, waiting for resp...\n", msg.RequestUID)

		select {
		case <-time.After(90 * time.Second):
			ctx.String(http.StatusGatewayTimeout, "Request timeout.")
			log.Printf("[WARN] Req timeout, ReqMsg: %v\n", msg)

		case resp := <-channel:
			log.Printf("[INFO]: Req has got the resp! ID is %v", msg.RequestUID)
			if err != nil {
				log.Printf("[Error]: return code is invalid")
				ctx.String(http.StatusInternalServerError, "HttpCode error")
				return
			}

			log.Printf("[INFO]: Resp received, body is %v, return type is %v",
				resp.Payload, resp.ContentType)
			ctx.Header("Content-Type", resp.ContentType)
			ctx.String(resp.StatusCode, resp.Payload)
		}
	}
}

func (rg *RuntimeGateway) SendMonitorInfo(action string) error {
	msg := message.MonitorMessage{
		InvokeTimeUnix: time.Now().Unix(),
		Action:         action,
	}

	log.Printf("[INFO]: Sending invoke info to monitor, action name is %v\n", action)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[Error]: Marshal monitor message failed\n")
		return err
	}

	err = rg.writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: options.MonitorTopic,
			Value: msgBytes,
		})
	if err != nil {
		log.Printf("[Error]: Write into monitor MQ failed\n")
		return err
	}

	return nil
}
