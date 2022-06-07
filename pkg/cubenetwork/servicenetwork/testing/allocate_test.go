package testing

import (
	"Cubernetes/pkg/cubenetwork/servicenetwork"
	"Cubernetes/pkg/object"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllocate(t *testing.T) {
	cia := servicenetwork.NewClusterIPAllocator()
	err := cia.Init()
	assert.NoError(t, err)

	service := &object.Service{
		TypeMeta:   object.TypeMeta{},
		ObjectMeta: object.ObjectMeta{},
		Spec: object.ServiceSpec{
			Selector:  nil,
			Ports:     nil,
			ClusterIP: "",
		},
		Status: nil,
	}

	// service, err = cia.AllocateClusterIP(service)
	// assert.Equal(t, "172.16.0.0", service.Spec.ClusterIP)
	// service, err = cia.AllocateClusterIP(service)
	// assert.Equal(t, "172.16.0.1", service.Spec.ClusterIP)
	// service.Spec.ClusterIP = "192.168.0.10"
	// service, err = cia.AllocateClusterIP(service)
	// assert.Equal(t, "172.16.0.2", service.Spec.ClusterIP)
	fmt.Println(service.Spec.ClusterIP)
}
