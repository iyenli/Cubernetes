package dns

import (
	"log"
	osexec "os/exec"
	"strings"
)

const (
	weaveName     = "weave"
	dnsAdd        = "dns-add"
	dnsRemove     = "dns-remove"
	defaultSuffix = ".weave.local"
)

func AddDNSEntry(containerName string, containerID string) error {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	// add default name
	var str strings.Builder
	str.WriteString(containerName)
	if !strings.HasSuffix(containerName, defaultSuffix) {
		str.WriteString(defaultSuffix)
	}

	containerName = str.String()
	containerName = strings.Trim(containerName, "\n")
	containerID = strings.Trim(containerID, "\n")

	cmd := osexec.Command(path, dnsAdd, containerID, "-h", containerName)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	if err != nil {
		return err
	}

	return nil
}

func DeleteDNSEntry(containerID string) error {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	cmd := osexec.Command(path, dnsRemove, containerID)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
