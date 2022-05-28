package actor_runtime

import (
	"Cubernetes/pkg/cubelet/cuberuntime/options"
	"Cubernetes/pkg/object"
	"log"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
)

const (
	defaultPauseImage = "docker/desktop-kubernetes-pause:3.7"
)

func (arm *actorRuntimeManager) startSandbox(actor *object.Actor) (string, string, error) {
	sandboxName := makeSandboxName(actor)
	sandboxConfig := &dockertypes.ContainerCreateConfig{
		Name: sandboxName,
		Config: &dockercontainer.Config{
			Image:  defaultPauseImage,
			Labels: newSandboxLabels(actor),
		},
		HostConfig: &dockercontainer.HostConfig{
			IpcMode:     dockercontainer.IpcMode("shareable"),
			DNS:         []string{options.WeaveDNSServer},
			DNSSearch:   []string{options.WeaveDNSSearchDomain},
			NetworkMode: options.WeaveNetwork,
		},
	}

	sandboxID, err := arm.dockerRuntime.CreateContainer(sandboxConfig)
	if err != nil {
		log.Printf("fail to create sandbox of actor %s\n", actor.Name)
		return "", "", err
	}

	err = arm.dockerRuntime.StartContainer(sandboxID)
	if err != nil {
		log.Printf("fail to start sandbox of actor %s\n", actor.Name)
		return "", "", err
	}

	return sandboxID, sandboxName, nil
}
