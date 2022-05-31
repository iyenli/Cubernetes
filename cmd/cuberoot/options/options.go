package options

const (
	BUILDDIR  = "./build/"
	ETCD      = "/usr/local/bin/etcd"
	CUBEROOT  = BUILDDIR + "cuberoot"
	APISERVER = BUILDDIR + "apiserver"
	CUBELET   = BUILDDIR + "cubelet"
	CUBEPROXY = BUILDDIR + "cubeproxy"
	MANAGER   = BUILDDIR + "manager"
	SCHEDULER = BUILDDIR + "scheduler"
	GATEWAY   = BUILDDIR + "gateway"
	BRAIN     = BUILDDIR + "brain"

	LOGDIR       = "/var/log/cubernetes/"
	APISERVERLOG = LOGDIR + "apiserver.log"
	CUBELETLOG   = LOGDIR + "cubelet.log"
	CUBEPROXYLOG = LOGDIR + "cubeproxy.log"
	MANAGERLOG   = LOGDIR + "manager.log"
	ETCDLOG      = LOGDIR + "etcd.log"
	SCHEDULERLOG = LOGDIR + "scheduler.log"
	GATEWAYLOG   = LOGDIR + "gateway.log"
	BRAINLOG     = LOGDIR + "brain.log"
)

const (
	GatewayImage = "yiyanleee/serverless-gateway:v1.5"
	Usage        = "usage"
	UsageLabel   = "ServerlessGatewayPod"
	GatewayIP    = "10.96.0.0"
)
