package options

const (
	BUILDDIR  = "./build/"
	CUBEROOT  = BUILDDIR + "cuberoot"
	APISERVER = BUILDDIR + "apiserver"
	CUBELET   = BUILDDIR + "cubelet"
	CUBEPROXY = BUILDDIR + "cubeproxy"
	MANAGER   = BUILDDIR + "manager"

	LOGDIR       = "/var/log/cubernetes"
	APISERVERLOG = LOGDIR + "apiserver.log"
	CUBELETLOG   = LOGDIR + "cubelet.log"
	CUBEPROXYLOG = LOGDIR + "cubeproxy.log"
	MANAGERLOG   = LOGDIR + "manager.log"
)
