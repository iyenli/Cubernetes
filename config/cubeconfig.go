package cubeconfig

import (
	"time"
)

const ETCDTimeout = time.Second
const ETCDAddr = "127.0.0.1:2379"

var APIServerIp = "127.0.0.1"

const APIServerPort = 8080
const HeartbeatPort = 8081

const (
	KindPod        = "Pod"
	KindService    = "Service"
	KindReplicaset = "ReplicaSet"
	KindNode       = "Node"
)

const DefaultApiVersion = "v1"

const CubeVersion = "v1.0"

const (
	//MetaDir  = "/var/log/cubernetes/"
	MetaDir  = "./cubernetes/"
	MetaFile = MetaDir + "meta"
)
