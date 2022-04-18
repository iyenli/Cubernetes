package cri

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	api "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"

	"Cubernetes/pkg/cubelet/cri/util"
)

type remoteRuntimeService struct {
	timeout       time.Duration
	runtimeClient runtimeapi.RuntimeServiceClient
}

func NewRemoteRuntimeService(endpoint string, connectionTimeout time.Duration) (api.RuntimeService, error) {
	log.Printf("Connecting to runtime service, endpoint = %s\n", endpoint)
	addr, dialer, err := util.GetAddressAndDialer(endpoint)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithContextDialer(dialer), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		log.Printf("Connect remote runtime failed, address = %s\n", addr)
		return nil, err
	}

	service := &remoteRuntimeService{
		timeout: connectionTimeout,
	}

	if err := service.establishConnection(conn); err != nil {
		return nil, err
	}

	return service, nil
}

/// impl RuntimeVersioner for remoteRuntimeService

func (r *remoteRuntimeService) Version(apiVersion string) (*runtimeapi.VersionResponse, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	return r.runtimeClient.Version(ctx, &runtimeapi.VersionRequest{
		Version: apiVersion,
	})
}

/// impl ContainerManager for remoteRuntimeService

func (r *remoteRuntimeService) CreateContainer(podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.CreateContainer(ctx, &runtimeapi.CreateContainerRequest{
		PodSandboxId:  podSandboxID,
		Config:        config,
		SandboxConfig: sandboxConfig,
	})

	if resp.ContainerId == "" || err != nil {
		log.Printf("CreateContainer in sandbox from runtime service failed, podSandBoxId = %s\n", podSandboxID)
		return "", err
	}

	return resp.ContainerId, nil
}

func (r *remoteRuntimeService) StartContainer(containerID string) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	_, err := r.runtimeClient.StartContainer(ctx, &runtimeapi.StartContainerRequest{
		ContainerId: containerID,
	})

	if err != nil {
		log.Printf("StartContainer from runtime service failed, containerID = %s\n", containerID)
		return err
	}

	return nil
}

func (r *remoteRuntimeService) StopContainer(containerID string, timeout int64) error {
	t := r.timeout + time.Duration(timeout)*time.Second
	ctx, cancel := getContextWithTimeout(t)
	defer cancel()

	_, err := r.runtimeClient.StopContainer(ctx, &runtimeapi.StopContainerRequest{
		ContainerId: containerID,
		Timeout:     timeout,
	})

	if err != nil {
		log.Printf("StopContainer from runtime service failed, containerID = %s\n", containerID)
		return err
	}

	return nil
}

func (r *remoteRuntimeService) RemoveContainer(containerID string) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	_, err := r.runtimeClient.RemoveContainer(ctx, &runtimeapi.RemoveContainerRequest{
		ContainerId: containerID,
	})

	if err != nil {
		log.Printf("RemoveContainer from runtime service failed, containerID = %s\n", containerID)
		return err
	}

	return nil
}

func (r *remoteRuntimeService) ListContainers(filter *runtimeapi.ContainerFilter) ([]*runtimeapi.Container, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.ListContainers(ctx, &runtimeapi.ListContainersRequest{
		Filter: filter,
	})

	if err != nil {
		log.Println("ListContainers from runtime service failed.")
		return nil, err
	}

	return resp.Containers, nil
}

func (r *remoteRuntimeService) ContainerStatus(containerID string, verbose bool) (*runtimeapi.ContainerStatusResponse, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.ContainerStatus(ctx, &runtimeapi.ContainerStatusRequest{
		ContainerId: containerID,
	})

	if err != nil {
		// No log since this will be called very often
		return nil, err
	}

	return resp, nil
}

func (r *remoteRuntimeService) UpdateContainerResources(containerID string, resources *runtimeapi.LinuxContainerResources) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	_, err := r.runtimeClient.UpdateContainerResources(ctx, &runtimeapi.UpdateContainerResourcesRequest{
		ContainerId: containerID,
		Linux:       resources,
	})

	if err != nil {
		log.Printf("UpdateContainerResources from runtime service failed, containerID = %s\n", containerID)
		return err
	}

	return nil
}

func (r *remoteRuntimeService) ExecSync(containerID string, cmd []string, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout != 0 {
		// Use timeout + default timeout (2 minutes) as timeout to leave some time for
		// the runtime to do cleanup.
		ctx, cancel = getContextWithTimeout(r.timeout + timeout)
	} else {
		ctx, cancel = getContextWithCancel()
	}
	defer cancel()

	resp, err := r.runtimeClient.ExecSync(ctx, &runtimeapi.ExecSyncRequest{
		ContainerId: containerID,
		Cmd:         cmd,
		Timeout:     int64(timeout.Seconds()),
	})

	if err != nil {
		log.Printf("ExecSync cmd from runtime service failed, containerID = %s, cmd = %s\n", containerID, cmd)
		return nil, nil, err
	}

	if resp.ExitCode != 0 {
		err = fmt.Errorf("command '%s' exited with %d: %s", strings.Join(cmd, " "), resp.ExitCode, resp.Stderr)
	}

	return resp.Stdout, resp.Stderr, err
}

func (r *remoteRuntimeService) Exec(req *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.Exec(ctx, req)
	if err != nil {
		log.Printf("Exec from runtime service failed, ContainerId = %s, cmd = %s\n", req.ContainerId, req.Cmd)
		return nil, err
	}

	if resp.Url == "" {
		err = errors.New("Url not set, Exec failed.")
		return nil, err
	}

	return resp, nil
}

func (r *remoteRuntimeService) Attach(req *runtimeapi.AttachRequest) (*runtimeapi.AttachResponse, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.Attach(ctx, req)
	if err != nil {
		log.Printf("Attach from runtime service failed, ContainerId = %s\n", req.ContainerId)
		return nil, err
	}

	if resp.Url == "" {
		err = errors.New("Url not set, Attach failed.")
		return nil, err
	}

	return resp, nil
}

func (r *remoteRuntimeService) ReopenContainerLog(containerID string) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	_, err := r.runtimeClient.ReopenContainerLog(ctx, &runtimeapi.ReopenContainerLogRequest{
		ContainerId: containerID,
	})

	if err != nil {
		log.Printf("ReopenContainerLog from runtime service failed, containerID = %s\n", containerID)
		return err
	}

	return nil
}

/// impl PodSandboxManager for remoteRuntimeService

func (r *remoteRuntimeService) RunPodSandbox(config *runtimeapi.PodSandboxConfig, runtimeHandler string) (string, error) {
	ctx, cancel := getContextWithTimeout(r.timeout * 2)
	defer cancel()

	resp, err := r.runtimeClient.RunPodSandbox(ctx, &runtimeapi.RunPodSandboxRequest{
		Config:         config,
		RuntimeHandler: runtimeHandler,
	})

	if err != nil {
		log.Println("RunPodSandbox from runtime service failed")
		return "", err
	}

	return resp.PodSandboxId, nil
}

func (r *remoteRuntimeService) StopPodSandbox(podSandBoxID string) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	_, err := r.runtimeClient.StopPodSandbox(ctx, &runtimeapi.StopPodSandboxRequest{
		PodSandboxId: podSandBoxID,
	})

	if err != nil {
		log.Printf("StopPodSandbox from runtime service failed, podSandBoxID = %s\n", podSandBoxID)
		return err
	}

	return nil
}

func (r *remoteRuntimeService) RemovePodSandbox(podSandBoxID string) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	_, err := r.runtimeClient.RemovePodSandbox(ctx, &runtimeapi.RemovePodSandboxRequest{
		PodSandboxId: podSandBoxID,
	})

	if err != nil {
		log.Printf("RemovePodSandbox from runtime service failed, podSandBoxID = %s\n", podSandBoxID)
		return err
	}

	return nil
}

func (r *remoteRuntimeService) PodSandboxStatus(podSandBoxID string, verbose bool) (*runtimeapi.PodSandboxStatusResponse, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.PodSandboxStatus(ctx, &runtimeapi.PodSandboxStatusRequest{
		PodSandboxId: podSandBoxID,
	})

	if err != nil {
		if verbose {
			log.Printf("PodSandboxStatus from runtime service failed, podSandBoxID = %s\n", podSandBoxID)
		}
		return nil, err
	}

	return resp, nil
}

func (r *remoteRuntimeService) ListPodSandbox(filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.ListPodSandbox(ctx, &runtimeapi.ListPodSandboxRequest{
		Filter: filter,
	})

	if err != nil {
		log.Println("ListPodSandbox with filter from runtime service failed")
		return nil, err
	}

	return resp.Items, nil
}

func (r *remoteRuntimeService) PortForward(req *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.PortForward(ctx, req)
	if err != nil {
		log.Printf("PortForward from runtime service failed, PodSandBoxId = %s\n", req.PodSandboxId)
		return nil, err
	}

	if resp.Url == "" {
		err = errors.New("Url not set, Exec failed.")
		return nil, err
	}

	return resp, nil
}

/// impl ContainerStatsManager for remoteRuntimeService

func (r *remoteRuntimeService) ContainerStats(containerID string) (*runtimeapi.ContainerStats, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.ContainerStats(ctx, &runtimeapi.ContainerStatsRequest{
		ContainerId: containerID,
	})

	if err != nil {
		log.Printf("ContainerStats from runtime service failed, containerId = %s\n", containerID)
		return nil, err
	}

	return resp.Stats, nil
}

func (r *remoteRuntimeService) ListContainerStats(filter *runtimeapi.ContainerStatsFilter) ([]*runtimeapi.ContainerStats, error) {
	ctx, cancel := getContextWithCancel()
	defer cancel()

	resp, err := r.runtimeClient.ListContainerStats(ctx, &runtimeapi.ListContainerStatsRequest{
		Filter: filter,
	})

	if err != nil {
		log.Println("ListContainerStats with filter from runtime service failed")
		return nil, err
	}

	return resp.Stats, nil
}

func (r *remoteRuntimeService) PodSandboxStats(podSandboxID string) (*runtimeapi.PodSandboxStats, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.PodSandboxStats(ctx, &runtimeapi.PodSandboxStatsRequest{
		PodSandboxId: podSandboxID,
	})

	if err != nil {
		log.Printf("PodSandboxStats from runtime service failed, podSandboxID = %s\n", podSandboxID)
		return nil, err
	}

	return resp.Stats, nil
}

func (r *remoteRuntimeService) ListPodSandboxStats(filter *runtimeapi.PodSandboxStatsFilter) ([]*runtimeapi.PodSandboxStats, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.ListPodSandboxStats(ctx, &runtimeapi.ListPodSandboxStatsRequest{
		Filter: filter,
	})

	if err != nil {
		log.Println("ListPodSandboxStats with filter from runtime service failed")
		return nil, err
	}

	return resp.Stats, nil
}

/// impl RuntimeService for remoteRuntimeService

func (r *remoteRuntimeService) UpdateRuntimeConfig(runtimeConfig *runtimeapi.RuntimeConfig) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	_, err := r.runtimeClient.UpdateRuntimeConfig(ctx, &runtimeapi.UpdateRuntimeConfigRequest{
		RuntimeConfig: runtimeConfig,
	})

	if err != nil {
		log.Println("UpdateRuntimeConfig from runtime service failed")
		return err
	}

	return err
}

func (r *remoteRuntimeService) Status(verbose bool) (*runtimeapi.StatusResponse, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.runtimeClient.Status(ctx, &runtimeapi.StatusRequest{})
	if err != nil {
		if verbose {
			log.Println("Status from runtime service failed")
		}
		return nil, err
	}

	return resp, nil
}

// establishConnection tries to connect to the remote runtime.
func (r *remoteRuntimeService) establishConnection(conn *grpc.ClientConn) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	r.runtimeClient = runtimeapi.NewRuntimeServiceClient(conn)

	if _, err := r.runtimeClient.Version(ctx, &runtimeapi.VersionRequest{}); err == nil {
		log.Println("Using CRI v1 runtime API")
	} else {
		return fmt.Errorf("unable to get runtime API version: %w", err)
	}

	return nil
}
