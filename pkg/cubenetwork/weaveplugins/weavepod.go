package weaveplugins

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/option"
	"errors"
	"log"
	"net"
	osexec "os/exec"
	"strings"
)

// AddPodToNetwork deprecated
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

// DeletePodFromNetwork deprecated
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

func GetPodIPByID(sandboxID string) (net.IP, error) {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return nil, err
	}

	log.Printf("[INFO]: Searching weave ip, sandbox id is %v", sandboxID)
	cmd := osexec.Command(path, "ps", sandboxID)
	byteOutput, err := cmd.CombinedOutput()
	output := strings.Trim(string(byteOutput), "\n")

	lines := strings.SplitAfter(output, "\n")
	if len(lines) != 1 {
		log.Printf("[Warn]: not one ip for contianer %v, ip number is %v", sandboxID, len(lines)-1)
		return nil, nil
	}

	cols := strings.Split(lines[0], " ")
	if len(cols) != 3 {
		log.Printf("[Warn]: not weave format, check it again, sandbox: %v", sandboxID)
		return nil, nil
	}

	ip, _, err := net.ParseCIDR(strings.Trim(cols[2], "\n"))
	if err != nil {
		log.Printf("[Warn]: not weave ip, check it again, sandbox: %v", sandboxID)
		return nil, nil
	}

	log.Printf("[INFO]: Searching weave ip, sandbox id is %v, ip is %v", sandboxID, ip.String())
	return ip, nil
}
