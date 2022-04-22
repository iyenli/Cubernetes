package cuberuntime

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	images "Cubernetes/pkg/cubelet/images"
	"log"

	criapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const (
	containerdRuntimeName = "containerd"
)

type cubeRuntimeManager struct {
	runtimeName string

	// wrapped image puller.
	imagePuller images.ImageManager

	// grpc service client
	runtimeService criapi.RuntimeService
	imageService   criapi.ImageManagerService
}

func (m *cubeRuntimeManager) PullImage(image cubecontainer.ImageSpec, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	// Pull without AuthConfig: not supported
	imageRef, err := m.imageService.PullImage(toRuntimeAPIImageSpec(image), nil, podSandboxConfig)
	if err != nil {
		log.Printf("fail to pull image #{image.Name}\n")
		return "", err
	}

	return imageRef, nil
}

func (m *cubeRuntimeManager) GetImageRef(image cubecontainer.ImageSpec) (string, error) {
	status, err := m.imageService.ImageStatus(toRuntimeAPIImageSpec(image))
	if err != nil {
		log.Printf("fail to get image #{image.Name} status\n")
		return "", err
	}

	if status == nil {
		return "", nil
	}

	return status.Id, nil
}

func (m *cubeRuntimeManager) ListImages() ([]cubecontainer.Image, error) {
	var images []cubecontainer.Image

	allImages, err := m.imageService.ListImages(nil)
	if err != nil {
		log.Printf("fail to list images\n")
		return images, err
	}

	for _, img := range allImages {
		images = append(images, cubecontainer.Image{
			ID:   img.Id,
			Size: int64(img.Size_),
			Spec: toCubeContainerImageSpec(img),
		})
	}

	return images, nil
}

func (m *cubeRuntimeManager) RemoveImage(image cubecontainer.ImageSpec) error {
	err := m.imageService.RemoveImage(&runtimeapi.ImageSpec{Image: image.Image})
	if err != nil {
		log.Printf("fail to remove image #{image.Name}\n")
		return err
	}

	return nil
}

type CubeRuntime interface {
}
