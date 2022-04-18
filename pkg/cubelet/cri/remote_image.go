package cri

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	api "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"

	"Cubernetes/pkg/cubelet/cri/util"
)

type remoteImageService struct {
	timeout     time.Duration
	imageClient runtimeapi.ImageServiceClient
}

func NewRemoteImageService(endpoint string, connectionTimeout time.Duration) (api.ImageManagerService, error) {
	log.Printf("Connecting to image service, endpoint = %s\n", endpoint)
	addr, dialer, err := util.GetAddressAndDialer(endpoint)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithContextDialer(dialer), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		log.Printf("Connect remote image failed, address = %s\n", addr)
		return nil, err
	}

	service := &remoteImageService{
		timeout: connectionTimeout,
	}

	if err := service.establishConnection(conn); err != nil {
		return nil, err
	}

	return service, nil
}

/// impl ImageManagerService for remoteImageService

func (r *remoteImageService) ListImages(filter *runtimeapi.ImageFilter) ([]*runtimeapi.Image, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.imageClient.ListImages(ctx, &runtimeapi.ListImagesRequest{
		Filter: filter,
	})

	if err != nil {
		log.Println("ListImages with filter from image service failed")
		return nil, err
	}

	return resp.Images, nil
}

func (r *remoteImageService) ImageStatus(image *runtimeapi.ImageSpec) (*runtimeapi.Image, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	resp, err := r.imageClient.ImageStatus(ctx, &runtimeapi.ImageStatusRequest{
		Image: image,
	})

	if err != nil {
		log.Printf("ImageStatus from image service failed, image = %s\n", image.Image)
		return nil, err
	}

	return resp.Image, nil
}

func (r *remoteImageService) PullImage(image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	ctx, cancel := getContextWithCancel()
	defer cancel()

	resp, err := r.imageClient.PullImage(ctx, &runtimeapi.PullImageRequest{
		Image:         image,
		Auth:          auth,
		SandboxConfig: podSandboxConfig,
	})

	if err != nil {
		log.Printf("PullImage from image service failed, image = %s\n", image.Image)
		return "", err
	}

	if resp.ImageRef == "" {
		return "", fmt.Errorf("PullImage failed: imageRef of image %q is not set", image.Image)
	}

	return resp.ImageRef, nil
}

func (r *remoteImageService) RemoveImage(image *runtimeapi.ImageSpec) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	_, err := r.imageClient.RemoveImage(ctx, &runtimeapi.RemoveImageRequest{
		Image: image,
	})

	if err != nil {
		log.Printf("RemoveImage from image service failed, image = %s\n", image.Image)
	}

	return nil
}

func (r *remoteImageService) ImageFsInfo() ([]*runtimeapi.FilesystemUsage, error) {
	ctx, cancel := getContextWithCancel()
	defer cancel()

	resp, err := r.imageClient.ImageFsInfo(ctx, &runtimeapi.ImageFsInfoRequest{})
	if err != nil {
		log.Println("ImageFsInfo from image service failed")
	}

	return resp.ImageFilesystems, nil
}

// establishConnection tries to connect to the remote runtime.
func (r *remoteImageService) establishConnection(conn *grpc.ClientConn) error {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	r.imageClient = runtimeapi.NewImageServiceClient(conn)

	if _, err := r.imageClient.ImageFsInfo(ctx, &runtimeapi.ImageFsInfoRequest{}); err == nil {
		log.Println("Using CRI v1 image API")
	} else {
		return fmt.Errorf("unable to get image API version: %w", err)
	}

	return nil
}
