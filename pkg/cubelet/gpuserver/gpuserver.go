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
		log.Println("[INFO]: Inspecting exist docker resources")
		jr.mutex.Lock()
		for job, container := range jr.jobMap {
			log.Println("[INFO]: Inspecting docker status and release exited docker")
			// TODO: Inspecting
			if container == "" {
				delete(jr.jobMap, job)
			}
		}
		jr.mutex.Unlock()

		time.Sleep(time.Second * 60)
	}
}
