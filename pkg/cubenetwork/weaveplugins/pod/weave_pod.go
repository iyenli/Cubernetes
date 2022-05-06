package pod

import (
	"errors"
	"log"
	"net"
	osexec "os/exec"
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
	output, err := cmd.CombinedOutput()

	ip := net.ParseIP(string(output))
	if ip == nil {
		log.Printf("Weave not return correct ip: %v", string(output))
		return nil, errors.New("weave not return correct ip")
	}

	return ip, nil
}

func DeletePodFromNetwork(sandboxID string) (net.IP, error) {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return nil, err
	}

	cmd := osexec.Command(path, detach, sandboxID)
	output, err := cmd.CombinedOutput()

	ip := net.ParseIP(string(output))
	if ip == nil {
		log.Printf("Weave not return correct ip: %v", string(output))
		return nil, errors.New("weave not return correct ip")
	}

	return ip, nil
}
