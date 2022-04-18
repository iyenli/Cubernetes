package cuberuntime

import (
	criapi "k8s.io/cri-api/pkg/apis"
)

const (
	containerdRuntimeName = "containerd"
)

type cubeRuntimeManager struct {
	runtimeName string

	// grpc service client
	runtimeService criapi.RuntimeService
	imageService   criapi.ImageManagerService
}

type CubeRuntime interface {
}
