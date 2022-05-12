package weaveplugins

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/option"
	"log"
	osexec "os/exec"
	"strings"
)

func AddDNSEntry(containerName string, containerID string) error {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	// add default name
	var str strings.Builder
	str.WriteString(containerName)
	if !strings.HasSuffix(containerName, option.DefaultSuffix) {
		str.WriteString(option.DefaultSuffix)
	}

	containerName = str.String()
	containerName = strings.Trim(containerName, "\n")
	containerID = strings.Trim(containerID, "\n")

	cmd := osexec.Command(path, option.DnsAdd, containerID, "-h", containerName)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	if err != nil {
		return err
	}

	return nil
}

func DeleteDNSEntry(containerID string) error {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	cmd := osexec.Command(path, option.DnsRemove, containerID)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
