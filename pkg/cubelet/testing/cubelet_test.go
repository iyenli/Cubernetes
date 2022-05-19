package main

import (
	"Cubernetes/pkg/cubelet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestSyncLoop(t *testing.T) {
	log.Println("Testing syncloop")
	cubeletInstance := cubelet.Cubelet{}

	go cubeletInstance.Run()

	time.Sleep(10 * time.Second)
	resp, err := http.Get("http://127.0.0.1:8080")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}
