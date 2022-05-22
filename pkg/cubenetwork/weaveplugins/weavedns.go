package weaveplugins

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/option"
	"log"
	osexec "os/exec"
	"strings"
)

func AddDNSEntry(hostname string, serviceIP string) error {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("[Error]: Weave Not found.")
		return err
	}

	hostname = GetDNSHost(hostname)
	serviceIP = strings.Trim(serviceIP, "\n")

	cmd := osexec.Command(path, option.DnsAdd, serviceIP, "-h", hostname)
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

func GetDNSHost(hostname string) string {
	// add default name
	var str strings.Builder
	str.WriteString(hostname)
	if !strings.HasSuffix(hostname, option.DefaultSuffix) {
		str.WriteString(option.DefaultSuffix)
	}

	hostname = str.String()
	hostname = strings.Trim(hostname, "\n")

	return hostname
}
