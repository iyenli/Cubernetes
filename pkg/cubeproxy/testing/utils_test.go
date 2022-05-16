package testing

import (
	"Cubernetes/pkg/cubeproxy/proxyruntime/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDNSCheck(t *testing.T) {
	str := "example.com/"
	ans, err := utils.CheckDNSHostName(str)
	assert.NoError(t, err)
	assert.Equal(t, "example.com", ans)
}
