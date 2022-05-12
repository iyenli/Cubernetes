package weaveplugins

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/option"
	"errors"
	"log"
	"net"
	osexec "os/exec"
	"strings"
)

func AddPodToNetwork(sandboxID string) (net.IP, error) {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return nil, err
	}

	log.Printf("New pod added to network, sandbox id is %v", sandboxID)
	cmd := osexec.Command(path, option.Attach, sandboxID)
	byteOutput, err := cmd.CombinedOutput()

	output := strings.Trim(string(byteOutput), "\n")
	output = strings.Trim(output, " ")

	log.Printf("New pod ip allocated: %v", output)
	ip := net.ParseIP(output)
	if ip == nil {
		log.Printf("Weave not return correct ip: %v", string(output))
		return nil, errors.New("weave not return correct ip")
	}

	return ip, nil
}

func DeletePodFromNetwork(sandboxID string) error {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	log.Printf("New pod deleteded from network, sandbox id is %v", sandboxID)
	cmd := osexec.Command(path, option.Detach, sandboxID)
	err = cmd.Run()

	return nil
}
