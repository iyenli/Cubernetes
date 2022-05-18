package testing

import (
	"Cubernetes/pkg/cubelet/gpuserver"
	"Cubernetes/pkg/object"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGPUServer(t *testing.T) {
	job := object.GpuJob{
		TypeMeta: object.TypeMeta{
			Kind:       "GpuJob",
			APIVersion: "v1",
		},
		ObjectMeta: object.ObjectMeta{
			Name: "test-gpu",
		},
		Spec: object.GpuJobSpec{
			Filename: "../example/gpujob/GpuJobTest.tar.gz",
		},
		Status: object.GpuJobStatus{},
	}

	jr := gpuserver.NewJobRuntime()
	err := jr.AddGPUJob(&job)
	assert.NoError(t, err)

}
