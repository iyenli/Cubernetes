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

	if err != nil {
		log.Println("fail to send http get request, err: ", err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http get response body, err: ", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP GET NOT OK, status code: %v, server response: %s\n", resp.StatusCode, string(body))
		return nil, errors.New(string(body))
	}

	return body, nil
}

func postRequest(url string, obj any) ([]byte, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		log.Println("fail to marshal object, err: ", err)
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Println("fail to send http post request, err: ", err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http post response body, err: ", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP POST NOT OK, status code: %v, server response: %s\n", resp.StatusCode, string(body))
		return nil, errors.New(string(body))
	}

	return body, nil
}

func putRequest(url string, obj any) ([]byte, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		log.Println("fail to marshal object, err: ", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(buf))
	if err != nil {
		log.Println("fail to create http put request, err: ", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("fail to send http put request, err: ", err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http put response body, err: ", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP PUT NOT OK, status code: %v, server response: %s\n", resp.StatusCode, string(body))
		return nil, errors.New(string(body))
	}

	return body, nil
}

func deleteRequest(url string) error {
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader("{}"))
	if err != nil {
		log.Println("fail to create http delete request, err: ", err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil || resp == nil {
		log.Println("fail to send http delete request, err: ", err)
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http delete response body, err: ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP DELETE NOT OK, status code: %v, server response: %s\n", resp.StatusCode, string(body))
		return errors.New(string(body))
	}

	return nil
}
