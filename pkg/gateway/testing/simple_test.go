package testing

import (
	"Cubernetes/pkg/gateway/types"
	"encoding/json"
	"github.com/stretchr/testify/assert"
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
	resp := types.MQMessageResponse{
		RequestUID:  "1d8fa4dd-216d-4c68-ab91-9c698144469a",
		ContentType: "application/json",
		StatusCode:  200,
		Payload:     "\"p\": \"a\"",
	}
	a, _ := json.Marshal(resp)
	log.Printf("%v", string(a))
}

func TestParseReturn(t *testing.T) {
	s := "{\"requestUID\": \"d32d4b93-fd6d-40a9-a8f7-ef68f9a66ee4\", \"statusCode\": 200, \"contentType\": \"text/plain\", \"payload\": \"60\"}"
	msg := types.MQMessageResponse{}
	err := json.Unmarshal([]byte(s), &msg)
	assert.NoError(t, err)

}
