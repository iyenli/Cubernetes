package pod

import (
	"errors"
	"log"
	"net"
	osexec "os/exec"
	"strings"
)

const (
	weaveName = "weave"
	attach    = "attach"
	detach    = "detach"
)

func AddPodToNetwork(sandboxID string) (net.IP, error) {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return nil, err
	}

	cmd := osexec.Command(path, attach, sandboxID)
	byteOutput, err := cmd.CombinedOutput()

	output := strings.Trim(string(byteOutput), "\n")
	output = strings.Trim(output, " ")

	ip := net.ParseIP(output)
	if ip == nil {
		log.Printf("Weave not return correct ip: %v", string(output))
		return nil, errors.New("weave not return correct ip")
	}

	return ip, nil
}

func DeletePodFromNetwork(sandboxID string) error {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	cmd := osexec.Command(path, detach, sandboxID)
	err = cmd.Run()

	return nil
}
