package testing

import (
	"Cubernetes/pkg/cubelet/dockershim"
	"Cubernetes/pkg/cubeproxy/utils"
	"Cubernetes/pkg/cubeproxy/utils/options"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNaming(t *testing.T) {
	str := "example.com"
	ans := utils.Hostname2NginxDockerName(str)
	assert.Equal(t, options.DockerNamePrefix+str, ans)
	ans = utils.NginxDockerName2Hostname(ans)
	assert.Equal(t, str, ans)
}

func TestPullImage(t *testing.T) {
	runtime, err := dockershim.NewDockerRuntime()
	assert.NoError(t, err)

	err = runtime.PullImage("nginx")
	assert.NoError(t, err)
}
