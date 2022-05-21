package cubeconfig

import (
	"time"
)

const ETCDTimeout = time.Second
const ETCDAddr = "127.0.0.1:2379"
const ServiceClusterIPRange = "172.16.0.0/16"

var APIServerIp = "127.0.0.1"

const APIServerPort = 8080
const HeartbeatPort = 8081

const DefaultApiVersion = "v1"

const CubeVersion = "v1.0"

const (
	JobFileDir = "/etc/cubernetes/apiserver/jobs/"
	MetaDir    = "/etc/cubernetes/cubernetes/"
	MetaFile   = MetaDir + "meta"
)
