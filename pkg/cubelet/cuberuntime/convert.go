package cuberuntime

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	"Cubernetes/pkg/cubelet/dockershim"
	"strings"
	"time"

	dockertypes "github.com/docker/docker/api/types"
)

func toSandboxStatus(dc *dockertypes.Container) *cubecontainer.SandboxStatus {
	status := &cubecontainer.SandboxStatus{
		Id:     dc.ID,
		Name:   dockershim.ParseSandboxName(dc.Names[0]),
		PodUID: dc.Labels[CubernetesPodUIDLabel],
		State:  toSandboxState(dc.Status),
	}

	return status
}

func toContainerStatus(dc *dockertypes.Container) *cubecontainer.ContainerStatus {
	status := &cubecontainer.ContainerStatus{
		ID: cubecontainer.ContainerID{
			Type: "docker",
			ID:   dc.ID,
		},
		Name:      dc.Labels[CubernetesContainerNameLabel],
		State:     toContainerState(dc.Status),
		CreatedAt: time.Unix(0, dc.Created),
		Image:     dc.Image,
		ImageID:   dc.ImageID,
	}

	return status
}

func toSandboxState(state string) cubecontainer.SandboxState {
	switch {
	case strings.HasPrefix(state, "Up"):
		return cubecontainer.SandboxStateReady
	default:
		return cubecontainer.SandboxStateNotReady
	}
}

func toContainerState(state string) cubecontainer.ContainerState {
	switch {
	case strings.HasPrefix(state, "Up"):
		return cubecontainer.ContainerStateRunning
	case strings.HasPrefix(state, "Exited"):
		return cubecontainer.ContainerStateExited
	case strings.HasPrefix(state, "Created"):
		return cubecontainer.ContainerStateCreated
	default:
		return cubecontainer.ContainerStateUnknown
	}
}

// State, ExitCode
func toContainerStateAndExitCode(jsonState *dockertypes.ContainerState) (cubecontainer.ContainerState, int) {
	if jsonState.Running {
		return cubecontainer.ContainerStateRunning, 0
	} else if jsonState.Status == "exited" {
		return cubecontainer.ContainerStateExited, jsonState.ExitCode
	} else if jsonState.Status == "created" {
		return cubecontainer.ContainerStateCreated, 0
	} else {
		return cubecontainer.ContainerStateUnknown, 0
	}
}
