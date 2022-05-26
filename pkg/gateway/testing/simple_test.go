package testing

import (
	"Cubernetes/pkg/gateway/types"
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestGin(t *testing.T) {
	channel := make(chan string)
	go func() {
		time.Sleep(10 * time.Second)
		channel <- "ppd"
	}()

	resp := <-channel
	log.Println(resp)
}

func TestJSON(t *testing.T) {
	resp := types.MQMessage{
		RequestUID:  "1d8fa4dd-216d-4c68-ab91-9c698144469a",
		ContentType: "application/json",
		StatusCode:  "200",
		Payload:     "\"p\": \"a\"",
	}
	a, _ := json.Marshal(resp)
	log.Printf("%v", string(a))
}
