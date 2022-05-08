package options

const (
	BUILDDIR  = "./build/"
	CUBEROOT  = BUILDDIR + "cuberoot"
	APISERVER = BUILDDIR + "apiserver"
	CUBELET   = BUILDDIR + "cubelet"
	CUBEPROXY = BUILDDIR + "cubeproxy"
	MANAGER   = BUILDDIR + "manager"
	ETCD      = "/usr/local/bin/etcd"

	LOGDIR       = "/var/log/cubernetes/"
	APISERVERLOG = LOGDIR + "apiserver.log"
	CUBELETLOG   = LOGDIR + "cubelet.log"
	CUBEPROXYLOG = LOGDIR + "cubeproxy.log"
	MANAGERLOG   = LOGDIR + "manager.log"
	ETCDLOG      = LOGDIR + "etcd.log"
)
