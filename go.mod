module Cubernetes

go 1.18

require k8s.io/cri-api v0.23.5

require (
	github.com/containerd/go-cni v1.1.4
	github.com/docker/docker v20.10.14+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/gin-gonic/gin v1.7.7
	github.com/google/uuid v1.3.0
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.7.1
	go.etcd.io/etcd/api/v3 v3.5.4
	go.etcd.io/etcd/client/v3 v3.5.4
	google.golang.org/grpc v1.43.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require gopkg.in/yaml.v2 v2.4.0

require (
	github.com/containerd/containerd v1.6.3 // indirect
	github.com/coreos/go-iptables v0.6.0
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/go-playground/validator/v10 v10.10.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/ugorji/go v1.2.7 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	google.golang.org/genproto v0.0.0-20220107163113-42d7afdf6368 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
