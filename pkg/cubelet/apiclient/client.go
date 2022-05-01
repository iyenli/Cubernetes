package apiclient

import (
	"Cubernetes/pkg/object"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type APIClient interface {
	PutPod(*object.Pod) error
}

func NewAPIClient(serverUrl string) (APIClient, error) {
	client := &aPIClient{serverUrl: serverUrl}

	// TODO: check url
	return client, nil
}

type aPIClient struct {
	serverUrl string
}

func (a *aPIClient) PutPod(pod *object.Pod) error {
	if pod.UID == "" {
		return fmt.Errorf("must put pod with UID")
	}

	payloadBytes, err := json.Marshal(pod)
	if err != nil {
		log.Printf("fail to marshall pod %s: %v", pod.Name, err)
		return err
	}

	// build http request
	url := fmt.Sprintf("%s/apis/pod/%s:%s", a.serverUrl, pod.Name, pod.UID)
	payload := strings.NewReader(string(payloadBytes))
	req, _ := http.NewRequest("PUT", url, payload)
	req.Header.Add("Content-Type", "application/json")

	resq, _ := http.DefaultClient.Do(req)
	if resq.StatusCode != 200 {
		return fmt.Errorf("fail to put pod to server, code = %d", resq.StatusCode)
	}

	return nil
}
