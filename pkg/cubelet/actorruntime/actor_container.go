package actor_runtime

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/cubelet/actorruntime/options"
	"Cubernetes/pkg/object"
	"log"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
)

func (arm *actorRuntimeManager) startContainer(actor *object.Actor, sandboxName string) error {
	mode := "container:" + sandboxName
	dirHostPath := arm.scriptManager.GetScriptDirPath(actor)

	config := &dockertypes.ContainerCreateConfig{
		Name: makeContainerName(actor),
		Config: &dockercontainer.Config{
			Image: options.ActorImageName,
			Cmd: append([]string{
				cubeconfig.APIServerIp + ":9092",
				actor.Spec.ActionName},
				actor.Spec.InvokeActions...),
			Labels: newContainerLabels(actor),
		},
		HostConfig: &dockercontainer.HostConfig{
			Binds: []string{
				strings.Join([]string{dirHostPath, options.ScriptDirMountPath}, ":")},
			NetworkMode: dockercontainer.NetworkMode(mode),
			IpcMode:     dockercontainer.IpcMode(mode),
			PidMode:     dockercontainer.PidMode(mode),
		},
	}

	containerID, err := arm.dockerRuntime.CreateContainer(config)
	if err != nil {
		log.Printf("fail to create container for actor %s\n", actor.Name)
		return err
	}

	err = arm.dockerRuntime.StartContainer(containerID)
	if err != nil {
		log.Printf("fail to start container for actor %s\n", actor.Name)
		return err
	}

	return nil
}
