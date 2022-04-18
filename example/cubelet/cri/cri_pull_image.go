package cri

// build: `go build -o build/pull_image Cubernetes/example/cubelet/cri`

import (
	"Cubernetes/pkg/cubelet/cri"
	"log"
	"os"
	"time"

	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const (
	containerdRuntimeEndpoint = "unix:///run/containerd/containerd.sock"
	connTimeout               = time.Second * 2
)

func main() {
	imageService, err := cri.NewRemoteImageService(containerdRuntimeEndpoint, connTimeout)
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	imageRef, err := imageService.PullImage(&v1.ImageSpec{Image: "busybox"}, nil, nil)
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	log.Println(imageRef)
	os.Exit(0)
}
