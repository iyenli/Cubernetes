package weaveplugins

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	osexec "os/exec"
)

const (
	weaveSource = "https://gitee.com/cderlee/weave/attach_files/1051080/download/weave"
	curl        = "curl"
	weavePath   = "/usr/local/bin/weave"
	chmod       = "chmod"
	exec        = "a+x"
)

// InstallWeave install weave in node
func InstallWeave() error {
	resp, err := http.Get(weaveSource)
	if err != nil {
		log.Printf("Http get weave failed: %v", err)
	}

	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(weavePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	write := bufio.NewWriter(file)
	_, err = write.Write(read)
	if err != nil {
		return err
	}

	err = write.Flush()
	if err != nil {
		return err
	}

	err = addExecute()
	if err != nil {
		log.Panicf("Install weave error: %s\n", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := file.Close()
		if err != nil {
			return
		}
		err = Body.Close()
		if err != nil {
			log.Println("Body closed failed")
		}
	}(resp.Body)

	return nil
}

func addExecute() error {
	path, err := osexec.LookPath(chmod)
	if err != nil {
		log.Panicln("Install chmod and use cubernetes.")
		return err
	}

	cmd := osexec.Command(path, exec, weavePath)
	err = cmd.Run()
	if err != nil {
		log.Panicf("Chmod weave error: %s\n", err)
		return err
	}

	return nil
}
