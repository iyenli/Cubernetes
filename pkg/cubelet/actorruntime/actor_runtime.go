package actor_runtime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/cubelet/dockershim"
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"Cubernetes/pkg/object"
	"log"
	"net"
	"strings"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type ActorRuntime interface {
	CreateActor(actor *object.Actor) error
	KillActor(UID string) error
	InspectActor(UID string) (object.ActorPhase, error)
	UpdateActionScript(actor *object.Actor) error
}

func NewActorRuntime() (ActorRuntime, error) {
	dockerRuntime, err := dockershim.NewDockerRuntime()
	if err != nil {
		log.Println("Fail to create docker client")
	}

	return &actorRuntimeManager{
		dockerRuntime: dockerRuntime,
		scriptManager: NewScriptManager(),
	}, nil
}

type actorRuntimeManager struct {
	dockerRuntime dockershim.DockerRuntime
	scriptManager ScriptManager
}

// assume all images exist
func (arm *actorRuntimeManager) CreateActor(actor *object.Actor) error {
	var ip net.IP

	if err := arm.scriptManager.EnsureScriptExist(actor); err != nil {
		log.Printf("fail to pull script of actor %s\n", actor.Name)
		return err
	}

	sandboxID, sandboxName, err := arm.startSandbox(actor)
	if err != nil {
		log.Printf("[Error] fail to start sandbox for actor %s\n", actor.Name)
		return err
	}

	if ip, err = weaveplugins.GetPodIPByID(sandboxID); err != nil || ip == nil {
		log.Printf("[Error]: add actor to weave network failed")
		return err
	} else {
		log.Printf("IP Allocated: %v", ip.String())
	}

	if err = arm.startContainer(actor, sandboxName); err != nil {
		log.Printf("[Error] fail to start container for actor %s\n", actor.Name)
		return err
	}

	actor.Status = &object.ActorStatus{
		Phase:           object.ActorRunning,
		IP:              ip,
		NodeUID:         actor.Status.NodeUID,
		LastUpdatedTime: time.Now(),
	}

	if _, err = crudobj.UpdateActor(*actor); err != nil {
		log.Printf("fail to update Actor %s status to apiserver\n", actor.Name)
		return err
	}

	return nil
}

func (arm *actorRuntimeManager) KillActor(UID string) error {
	// Only need to remove sandbox
	filter := dockertypes.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", buildLabelSelector(ContainerTypeLabel, ContainerTypeSandbox)),
			filters.Arg("label", buildLabelSelector(ActorUIDLabel, UID)),
		),
	}

	sandboxes, err := arm.dockerRuntime.ListContainers(filter)
	if err != nil {
		log.Printf("fail to list actor sandbox %s: %v\n", UID, err)
		return err
	}

	for _, sandbox := range sandboxes {
		if err := arm.dockerRuntime.StopContainer(sandbox.ID); err != nil {
			log.Printf("[Error]: fail to stop sandbox %s: %v\n", sandbox.ID, err)
			return err
		}

		if err := arm.dockerRuntime.RemoveContainer(sandbox.ID, false); err != nil {
			log.Printf("[Error]: fail to remove sandbox %s: %v\n", sandbox.ID, err)
			return err
		}
	}

	return nil
}

func (arm *actorRuntimeManager) UpdateActionScript(actor *object.Actor) error {
	// pull script if necessary
	if err := arm.scriptManager.EnsureScriptExist(actor); err != nil {
		log.Printf("[Error]: fail to update script of  %s: %v\n", actor.Spec.ActionName, err)
		return err
	}

	// get containerID
	containerID, err := arm.getActorContainerID(actor.UID)
	if err != nil {
		return err
	}

	// send SIGINT to actor container,
	// this will trigger wrapper.py to reload action module,
	// which updates the `action` function in time
	if err = arm.dockerRuntime.SignalContainer(containerID, "SIGINT"); err != nil {
		return err
	}

	log.Printf("script of actor %s updated!\n", actor.Name)

	return nil
}

func (arm *actorRuntimeManager) InspectActor(UID string) (object.ActorPhase, error) {
	filter := dockertypes.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", buildLabelSelector(ActorUIDLabel, UID)),
		),
	}

	// list both sandbox and container
	containers, err := arm.dockerRuntime.ListContainers(filter)
	if err != nil {
		log.Printf("fail to list actor sandbox and container %s: %v\n", UID, err)
		return object.ActorUnknown, err
	}

	running := 0
	for _, container := range containers {
		if strings.HasPrefix(container.Status, "Up") {
			running += 1
		}
	}

	if running == 2 {
		return object.ActorRunning, nil
	} else if running < 2 {
		return object.ActorFailed, nil
	} else {
		return object.ActorUnknown, nil
	}
}
