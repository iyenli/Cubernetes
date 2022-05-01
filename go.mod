module Cubernetes

go 1.18

require k8s.io/cri-api v0.23.5

require (
	github.com/containerd/go-cni v1.1.4
	github.com/coreos/etcd v3.3.27+incompatible
	github.com/gin-gonic/gin v1.7.7
	github.com/google/uuid v1.2.0
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.7.1
	go.etcd.io/etcd v3.3.27+incompatible
	google.golang.org/grpc v1.43.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/containerd/containerd v1.6.3 // indirect
	github.com/coreos/go-iptables v0.6.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.14+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/go-playground/validator/v10 v10.10.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/ugorji/go v1.2.7 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	google.golang.org/genproto v0.0.0-20220107163113-42d7afdf6368 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

replace google.golang.org/grpc v1.40.0 => google.golang.org/grpc v1.26.0
