package images

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	object "Cubernetes/pkg/object"
	"fmt"
	"log"
	"strings"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type imageManager struct {
	imageService cubecontainer.ImageService
}

type ImageManager = *imageManager

func NewImageManager(imageService cubecontainer.ImageService) ImageManager {
	return &imageManager{imageService}
}

func (m *imageManager) EnsureImageExists(pod *object.Pod, container *object.Container, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	log.Printf("Pulling image: #{pod.NameSpace}/#{pod.Name}/#{container.Image}\n")

	image := applyDefaultImageTag(container.Image)
	spec := cubecontainer.ImageSpec{Image: image}

	imageRef, err := m.imageService.GetImageRef(spec)
	if err != nil {
		return "", fmt.Errorf("failed to inspect image %q: %v", container.Image, err)
	}

	if present := imageRef != ""; present {
		log.Printf("Container image #{image} already present on machine\n")
		return imageRef, nil
	}

	imageRef, err = m.imageService.PullImage(spec, podSandboxConfig)
	if err != nil {
		log.Printf("Failed to pull image #{image}: #{err.Error()}\n")
		return "", err
	}

	return imageRef, nil
}

// if images doesn't contain any tag, apply "latest"
func applyDefaultImageTag(image string) string {
	if !strings.Contains(image, ":") {
		return image + ":latest"
	} else {
		return image
	}
}
