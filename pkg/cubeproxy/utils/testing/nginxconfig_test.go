package testing

import (
	"Cubernetes/pkg/cubeproxy/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNginxConfig(t *testing.T) {
	err := utils.CreateNginxConfig("test.test",
		[]string{"/path"}, []string{"10.32.0.2"}, []string{"8080"})
	assert.NoError(t, err)
}

func TestWriteFiles(t *testing.T) {
	var s = "123"
	err := utils.PrepareNginxFile("hello", &s)
	assert.NoError(t, err)
}
