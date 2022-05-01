package crudobj

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func getRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		log.Println("fail to send http get request")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http get response body")
		return nil, err
	}
	return body, nil
}

func postRequest(url string, obj any) ([]byte, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		log.Println("fail to marshal object")
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		log.Println("fail to send http post request")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http post response body")
		return nil, err
	}
	return body, nil
}

func putRequest(url string, obj any) ([]byte, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		log.Println("fail to marshal object")
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(buf))
	if err != nil {
		log.Println("fail to create http put request")
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		log.Println("fail to send http put request")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http put response body")
		return nil, err
	}
	return body, nil
}

func deleteRequest(url string) error {
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader("{}"))
	if err != nil {
		log.Println("fail to create http delete request")
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		log.Println("fail to send http delete request")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http delete response body")
		return err
	}

	if string(body) != "\"deleted\"" {
		return errors.New("fail to delete the obj")
	}

	return nil
}
