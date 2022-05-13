package cuberuntime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	cubecontainer "Cubernetes/pkg/cubelet/container"
	dockershim "Cubernetes/pkg/cubelet/dockershim"
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	object "Cubernetes/pkg/object"
	"fmt"
	"log"
	"time"
)

const (
	containerdRuntimeName     = "containerd"
	podLogsRootDirectory      = "/var/log/pods"
	containerdRuntimeEndpoint = "unix:///run/containerd/containerd.sock"
	remoteConnectTimeout      = time.Second * 2
)

type cubeRuntimeManager struct {
	runtimeName string

	dockerRuntime dockershim.DockerRuntime
}

type podActions struct {
	KillPod       bool
	CreateSandbox bool
	// old sandbox id, kill if we need to kill old pod
	SandboxID string
	// index of containers in podSpec.Containers to start
	ContainersToStart []int
	// UID of containers to kill
	ContainersToKill []string
}

type CubeRuntime interface {
	cubecontainer.Runtime
}

func (m *cubeRuntimeManager) SyncPod(pod *object.Pod, podStatus *cubecontainer.PodStatus) error {

	// Compute sandbox and container changes.
	podContainerChanges := m.computePodActions(pod, podStatus)

	log.Printf("\ncreate sandbox: %t\ncreate container: %v\n\n", podContainerChanges.CreateSandbox, podContainerChanges.ContainersToStart)

	removeContainer := true
	// Kill the pod if sandbox changed
	if podContainerChanges.KillPod {
		if err := m.killPodByStatus(podStatus, removeContainer); err != nil {
			log.Printf("fail to kill pod %s: %v\n", pod.Name, err)
			return err
		}
	} else {
		// kill some containers
		for _, uid := range podContainerChanges.ContainersToKill {
			if err := m.dockerRuntime.StopContainer(uid); err != nil {
				log.Printf("fail to kill container uid %s: %v\n", uid, err)
				return err
			}
		}
	}

	// Create sandbox if necessary
	podSandboxID := podContainerChanges.SandboxID
	podSandboxName := dockershim.MakeSandboxName(pod)
	if podContainerChanges.CreateSandbox {
		var err error

		if podSandboxName, podSandboxID, err = m.createPodSandbox(pod); err != nil {
			return err
		}
		log.Printf("create sandbox %s for pod %s\n", podSandboxID, pod.Name)

		// Update sandbox to initnetwork
		newSandboxStatuses, _ := m.getSandboxStatusesByPodUID(pod.UID)
		podStatus.UpdateSandboxStatuses(newSandboxStatuses)

		ip, err := weaveplugins.AddPodToNetwork(podSandboxID)
		if err != nil || ip == nil {
			log.Printf("[Error]: add pod to weave network failed")
			return err
		}
		log.Printf("IP Allocated: %v", ip.String())
		//network.InitNetwork(network.ProbeNetworkPlugins("", ""), podStatus)

		podStatus.PodNetWork.IP = ip
	}

	// Create containers
	for _, idx := range podContainerChanges.ContainersToStart {
		msg, err := m.startContainer(&pod.Spec.Containers[idx], pod, podSandboxName)
		if err != nil {
			log.Printf("fail to start container %s: %s\n", pod.Spec.Containers[idx].Name, msg)
			return err
		}
		log.Printf("start container %s\n", pod.Spec.Containers[idx].Name)
	}

	apiPodStatus, err := m.InspectPod(pod)
	if err != nil {
		log.Printf("fail to get pod status %s: %v\n", pod.UID, err)
		return err
	}

	apiPodStatus.IP = podStatus.PodNetWork.IP
	if pod.Status != nil {
		apiPodStatus.PodUID = pod.Status.PodUID
	}

	_, err = crudobj.UpdatePodStatus(pod.UID, *apiPodStatus)
	if err != nil {
		log.Printf("fail to update Pod %s status to apiserver\n", pod.Name)
		return err
	} else {
		log.Printf("update Pod %s status by SyncPod\n", pod.Name)
	}

	return nil
}

// FIXME: only compute container changes by its name now
func (m *cubeRuntimeManager) computePodActions(pod *object.Pod, podStatus *cubecontainer.PodStatus) podActions {
	createPodSandbox, sandboxID := m.podSandboxChanged(pod, podStatus)
	changes := podActions{
		KillPod:           createPodSandbox,
		CreateSandbox:     createPodSandbox,
		SandboxID:         sandboxID,
		ContainersToStart: []int{},
		ContainersToKill:  []string{},
	}

	// create sandbox need to (re-)create all containers
	if createPodSandbox {
		var containersToStart []int
		var containersToKill []string

		for idx := range pod.Spec.Containers {
			// TODO: RestartPolicy == OnFailure && ExitSucceeded => no need to start
			containersToStart = append(containersToStart, idx)
		}

		for _, oldContainer := range podStatus.ContainerStatuses {
			// kill all old containers
			containersToKill = append(containersToKill, oldContainer.ID.ID)
		}

		if len(containersToStart) == 0 {
			// nothing to create, so don't create sandbox
			changes.CreateSandbox = false
		}

		changes.ContainersToStart = containersToStart
		changes.ContainersToKill = containersToKill
		return changes
	}

	var remain []string
	for _, oldContainer := range podStatus.ContainerStatuses {
		remain = append(remain, oldContainer.ID.ID)
	}

	for idx, container := range pod.Spec.Containers {
		containerStatus := podStatus.FindContainerStatusByName(container.Name)

		if containerStatus == nil || containerStatus.State != cubecontainer.ContainerStateRunning {
			// container not exist or container not running => simply restart
			// assume no name change => no spec change for easy impl
			changes.ContainersToStart = append(changes.ContainersToStart, idx)
			if containerStatus != nil /* just kill all old containers now */ {
				changes.ContainersToKill = append(changes.ContainersToKill, containerStatus.ID.ID)
			}
		} else {
			// container:name no change: keep the old container
			var keep int
			for idx := range remain {
				if remain[idx] == containerStatus.ID.ID {
					keep = idx
					break
				}
			}
			remain = append(remain[:keep], remain[keep+1:]...)
		}
	}

	// kill not mentioned containers
	changes.ContainersToKill = append(changes.ContainersToKill, remain...)

	return changes
}

// podSandboxChanged checks whether the spec of the pod is changed and returns
// (changed, original sandboxID if exist).
func (m *cubeRuntimeManager) podSandboxChanged(pod *object.Pod, podStatus *cubecontainer.PodStatus) (bool, string) {
	if len(podStatus.SandboxStatuses) == 0 {
		log.Printf("no sandbox for pod %s can be found. Need to start a new one.", pod.Name)
		// This branch should return
		return true, ""
	}

	sandboxStatus := podStatus.SandboxStatuses[0]
	if sandboxStatus.State != cubecontainer.SandboxStateReady {
		// No ready sandbox for pod can be found. Need to start a new one.
		return true, sandboxStatus.Id
	}

	// Needs to create a new sandbox when the sandbox does not have an IP address.
	if sandboxStatus.Ip == "" {
		// Sandbox for pod has no IP address. Need to start a new one.
		return true, sandboxStatus.Id
	}

	// sandbox unchange and still running
	return false, sandboxStatus.Id
}

func (m *cubeRuntimeManager) KillPod(UID string) error {
	log.Printf("Kill pod %s\n", UID)
	podStatus, err := m.getPodStatusByUID(UID)
	if err != nil {
		log.Printf("fail to get podStatus by UID %s\n", UID)
		return err
	}
	// for debug only
	removeContainer := true

	return m.killPodByStatus(podStatus, removeContainer)
}

func (m *cubeRuntimeManager) killPodByStatus(status *cubecontainer.PodStatus, remove bool) error {
	m.killPodContainers(status, remove)

	// kill pod sandbox
	for _, sandbox := range status.SandboxStatuses {
		log.Printf("start to kill sandbox %s\n", sandbox.Id)
		err := weaveplugins.DeletePodFromNetwork(sandbox.Id)
		if err != nil {
			return err
		}

		if err := m.dockerRuntime.StopContainer(sandbox.Id); err != nil {
			log.Printf("fail to stop sandbox %s: %v\n", sandbox.Id, err)
			return err
		}

		if remove {
			if err := m.dockerRuntime.RemoveContainer(sandbox.Id, false); err != nil {
				log.Printf("fail to remove sandbox %s: %v\n", sandbox.Id, err)
				return err
			}
		}

		//network.ReleaseNetwork(network.ProbeNetworkPlugins("", ""), status)
	}

	return nil
}

func (m *cubeRuntimeManager) GetPodStatus(UID string) (*cubecontainer.PodStatus, error) {
	return m.getPodStatusByUID(UID)
}

func (m *cubeRuntimeManager) InspectPod(pod *object.Pod) (*object.PodStatus, error) {
	containerStatuses, err := m.getContainerStatusesByPodUID(pod.UID)
	if err != nil {
		return nil, err
	}

	sandboxStatus, err := m.getSandboxStatusesByPodUID(pod.UID)
	if err != nil {
		return nil, err
	}

	if len(sandboxStatus) == 0 {
		return nil, fmt.Errorf("no sandbox for pod %s found", pod.Name)
	}

	// TODO: get sandbox IP from Network Plugin (Weave)
	var sandboxIP []byte
	podPhase := cubecontainer.ComputePodPhase(containerStatuses, sandboxStatus[0], &pod.Spec)

	usage := &object.ResourceUsage{}
	for _, status := range containerStatuses {
		usage.ActualCPUUsage += status.ResourceUsage.CPUUsage
		usage.ActualMemoryUsage += status.ResourceUsage.MemoryUsage
	}

	return &object.PodStatus{
		IP:                  sandboxIP,
		Phase:               podPhase,
		ActualResourceUsage: usage,
		LastUpdateTime:      time.Now(),
	}, nil
}

func (m *cubeRuntimeManager) ListPodsUID() ([]string, error) {
	return m.getAllPodsUID()
}

func (c *cubeRuntimeManager) getPodStatusByUID(UID string) (*cubecontainer.PodStatus, error) {
	containerStatuses, err := c.getContainerStatusesByPodUID(UID)
	if err != nil {
		return nil, err
	}

	sandboxStatuses, err := c.getSandboxStatusesByPodUID(UID)
	if err != nil {
		return nil, err
	}

	podName := ""
	if len(sandboxStatuses) > 0 {
		podName = sandboxStatuses[0].Name
	}

	return &cubecontainer.PodStatus{
		UID:              UID,
		Name:             podName,
		NetworkNamespace: "/var/run/netns/default",
		// Update PodNetworkStatus?
		ContainerStatuses: containerStatuses,
		SandboxStatuses:   sandboxStatuses,
	}, nil
}

func NewCubeRuntimeManager() (CubeRuntime, error) {
	dockerRuntime, err := dockershim.NewDockerRuntime()
	if err != nil {
		log.Println("Fail to create docker client")
	}

	cm := &cubeRuntimeManager{
		dockerRuntime: dockerRuntime,
		runtimeName:   containerdRuntimeName,
	}

	return cm, nil
}

func (c *cubeRuntimeManager) Close() {
	c.dockerRuntime.CloseConnection()
}
