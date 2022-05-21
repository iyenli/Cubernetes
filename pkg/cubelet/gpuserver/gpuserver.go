package gpuserver

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/cubelet/dockershim"
	"Cubernetes/pkg/cubelet/gpuserver/options"
	"Cubernetes/pkg/object"
	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"log"
	"sync"
	"time"
)

type JobRuntime struct {
	// Job UID -> Container ID
	jobMap         map[string]string
	dockerInstance dockershim.DockerRuntime

	mutex sync.Mutex
}

func NewJobRuntime() JobRuntime {
	dockerInstance, err := dockershim.NewDockerRuntime()
	if err != nil {
		log.Println("[Error]: Init docker runtime error")
	}

	return JobRuntime{
		jobMap:         make(map[string]string),
		dockerInstance: dockerInstance,
		mutex:          sync.Mutex{},
	}
}

func (jr *JobRuntime) AddGPUJob(job *object.GpuJob) error {
	jr.mutex.Lock()
	defer jr.mutex.Unlock()

	dockerName := GetJobDockerName(job.UID)
	log.Println("[INFO]: Creating docker, name:", dockerName)
	log.Println("[INFO]: Pulling Image", options.GpuServerImageName)

	err := jr.dockerInstance.PullImage(options.GpuServerImageName)
	if err != nil {
		log.Printf("[INFO]: Pull image failed\n")
		return err
	}

	config := &dockertypes.ContainerCreateConfig{
		Name: dockerName,
		Config: &dockercontainer.Config{
			Image: options.GpuServerImageName,
			Cmd:   strslice.StrSlice{job.UID, cubeconfig.APIServerIp},
		},
	}

	containerID, err := jr.dockerInstance.CreateContainer(config)
	if err != nil {
		log.Printf("[Error]: fail to create container %v\n", dockerName)
		return err
	}

	err = jr.dockerInstance.StartContainer(containerID)
	if err != nil {
		log.Printf("[Error]: fail to start container %v\n", dockerName)
		return err
	}

	jr.jobMap[job.UID] = containerID
	return nil
}

func (jr *JobRuntime) ReleaseContainerResource() {
	for {
		jr.mutex.Lock()
		log.Println("[INFO]: Inspecting docker status and release exited docker every 2 minutes...")
		for job, container := range jr.jobMap {
			if container == "" {
				delete(jr.jobMap, job)
			}
			log.Printf("[INFO]: Clearing job %v, corresponding containerID is %v",
				job, container)

			inspectContainer, err := jr.dockerInstance.InspectContainer(container)
			if err != nil {
				log.Printf("[INFO]: Inspect container ID %v failed", container)
				continue
			}

			if inspectContainer.State.Status == "exited" {
				log.Printf("[INFO]: container ID %v exited, it would be removed soon...", container)
				// Not actually do by now:)
			}
		}
		jr.mutex.Unlock()

		time.Sleep(time.Second * 120)
	}
}
