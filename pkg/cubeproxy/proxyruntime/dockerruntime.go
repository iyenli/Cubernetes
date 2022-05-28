package proxyruntime

import (
	"Cubernetes/pkg/cubelet/cuberuntime/options"
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"Cubernetes/pkg/cubeproxy/utils"
	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"log"
	"strings"
)

const (
	NginxImageName  = "nginx"
	NginxConfigPath = "/etc/nginx/"
)

// StartDNSNginxDocker For every DNS, Start a Nginx Docker
func (pr *ProxyRuntime) StartDNSNginxDocker(host string, paths, serviceIPs, ports []string) (string, error) {
	dockerName := utils.Hostname2NginxDockerName(host)
	err := utils.CreateNginxConfig(host, paths, serviceIPs, ports)
	if err != nil {
		log.Println("[Error]: Create Nginx config file failed")
		return "", err
	}

	log.Println("[INFO]: Creating docker, name:", dockerName)
	log.Println("[INFO]: Pulling Image", NginxImageName)

	err = pr.DockerInstance.PullImage("nginx")
	if err != nil {
		log.Printf("[INFO]: Pull nginx failed\n")
		return "", err
	}

	// prepare config bind
	volumeBinds := make([]string, 0)
	volumeBinds = append(volumeBinds,
		strings.Join([]string{utils.NginxConfigFileLocation(host), NginxConfigPath}, ":"),
	)

	config := &dockertypes.ContainerCreateConfig{
		Name: dockerName,
		Config: &dockercontainer.Config{
			Image:    NginxImageName,
			Hostname: weaveplugins.GetDNSHost(host),
		},
		HostConfig: &dockercontainer.HostConfig{
			Binds:       volumeBinds,
			DNS:         []string{options.WeaveDNSServer},
			DNSSearch:   []string{options.WeaveDNSSearchDomain},
			NetworkMode: options.WeaveNetwork,
		},
	}

	containerID, err := pr.DockerInstance.CreateContainer(config)
	if err != nil {
		log.Printf("[Error]: fail to create container %v\n", dockerName)
		return "", err
	}

	err = pr.DockerInstance.StartContainer(containerID)
	if err != nil {
		log.Printf("[Error]: fail to start container %v\n", dockerName)
		return "", err
	}

	return containerID, nil
}
