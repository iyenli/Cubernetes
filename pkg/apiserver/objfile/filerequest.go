package objfile

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

func getFile(url string, filename string) error {
	resp, err := http.Get(url)

	if err != nil {
		log.Println("fail to send http get request, err: ", err)
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http get response body, err: ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP GET NOT OK when getting file, status code: %v, server response: %s\n", resp.StatusCode, string(body))
		return errors.New(string(body))
	}

	err = ioutil.WriteFile(filename, body, 0777)
	if err != nil {
		log.Printf("fail to write into file: %s, err: %v\n", filename, err)
		return err
	}

	return nil
}

func postFile(url string, filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("fail to read from file: %s, err: %v\n", filename, err)
		return err
	}

	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Println("fail to create buffer writer, err: ", err)
		return err
	}

	_, err = fileWriter.Write(buf)
	if err != nil {
		log.Println("fail to write buffer, err: ", err)
		return err
	}

	err = writer.Close()
	if err != nil {
		log.Println("fail to close buffer writer, err: ", err)
		return err
	}

	resp, err := http.Post(url, writer.FormDataContentType(), &body)

	if err != nil {
		log.Println("fail to send http post request, err: ", err)
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http post response body, err: ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP POST NOT OK when posting file, status code: %v, server response: %s\n", resp.StatusCode, string(respBody))
		return errors.New(string(respBody))
	}

	return nil
}
