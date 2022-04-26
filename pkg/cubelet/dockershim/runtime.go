package dockershim

import (
	"context"
	"fmt"
	"log"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockerapi "github.com/docker/docker/client"
)

type DockerRuntime interface {
	// Container Service
	CreateContainer(config *dockertypes.ContainerCreateConfig) (string, error)
	StopContainer(containerID string) error

	// Image Service
	PullImage(imageName string) error
	RemoveImage(imageName string) error
	ListImages(all bool) ([]*dockertypes.ImageSummary, error)

	// Closer
	CloseConnection()
}

func NewDockerRuntime() (DockerRuntime, error) {
	client, err := dockerapi.NewClientWithOpts(dockerapi.FromEnv)
	if err != nil {
		log.Println("fail to connect docker from env")
		return nil, err
	}

	cubeDockerClient := &dockerClient{
		client:            client,
		timeout:           time.Minute * 2,
		imagePullDeadline: time.Minute,
	}

	return cubeDockerClient, nil
}

type dockerClient struct {
	timeout           time.Duration
	imagePullDeadline time.Duration
	client            *dockerapi.Client
}

func (c *dockerClient) CreateContainer(config *dockertypes.ContainerCreateConfig) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.ContainerCreate(ctx,
		config.Config,
		config.HostConfig,
		config.NetworkingConfig,
		config.Platform,
		config.Name)
	if err != nil {
		log.Printf("fail to create container %s : %v\n", config.Name, err)
		return "", err
	}

	if len(resp.Warnings) > 0 {
		log.Print("[Waring] ", resp.Warnings)
	}
	return resp.ID, nil
}

func (c *dockerClient) StopContainer(containerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	if err := c.client.ContainerStop(ctx, containerID, &c.timeout); err != nil {
		log.Printf("fail to stop container %s : %v\n", containerID, err)
		return err
	}

	return nil
}

func (c *dockerClient) PullImage(imageName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout*2)
	defer cancel()

	out, err := c.client.ImagePull(ctx, imageName, dockertypes.ImagePullOptions{})
	if err != nil {
		log.Printf("fail to pull image %s : %v\n", imageName, err)
		return err
	}

	defer out.Close()
	return nil
}

func (c *dockerClient) ListImages(all bool) ([]*dockertypes.ImageSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	images, err := c.client.ImageList(ctx, dockertypes.ImageListOptions{All: all})
	if err != nil {
		log.Printf("fail to list images: %v\n", err)
		return nil, err
	}

	imageRefs := []*dockertypes.ImageSummary{}
	for _, image := range images {
		imageRefs = append(imageRefs, &image)
	}

	return imageRefs, nil
}

func (c *dockerClient) RemoveImage(imageName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	images, err := c.client.ImageList(ctx, dockertypes.ImageListOptions{All: true})
	if err != nil {
		log.Printf("fail to list images: %v\n", err)
		return err
	}

	imageID := ""
	for _, i := range images {
		if i.RepoTags[0] == imageName {
			imageID = i.ID
		}
	}
	if imageID == "" {
		log.Printf("fail to find image: %s\n", imageName)
		return fmt.Errorf("fail to find image: %s", imageName)
	}

	resps, err := c.client.ImageRemove(ctx, imageID, dockertypes.ImageRemoveOptions{})
	if err != nil {
		log.Printf("fail to remove image: %v\n", err)
		return err
	}

	for _, resp := range resps {
		log.Println(resp.Deleted, resp.Untagged)
	}
	return nil
}

func (c *dockerClient) CloseConnection() {
	c.client.Close()
}
