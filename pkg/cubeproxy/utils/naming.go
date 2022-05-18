package utils

import (
	"Cubernetes/pkg/cubeproxy/utils/options"
	"log"
	"strings"
)

func Hostname2NginxDockerName(hostname string) string {
	log.Println("[INFO]: Docker name:", options.DockerNamePrefix+hostname)
	return options.DockerNamePrefix + hostname
}

func NginxDockerName2Hostname(hostname string) string {
	if !strings.HasPrefix(hostname, options.DockerNamePrefix) {
		log.Println("[Warn]: illegal name sent into converter:", hostname)
	}

	hostname = hostname[len(options.DockerNamePrefix):]
	log.Println("[INFO]: host name:", hostname)
	return hostname
}

func NginxConfigFileLocation(hostname string) string {
	str := options.NginxFile + hostname + "/"
	return str
}
