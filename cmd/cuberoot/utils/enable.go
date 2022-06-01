package utils

import (
	"Cubernetes/cmd/cuberoot/options"
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/localstorage"
	"log"
)

func EnableServerlessGateway(meta *localstorage.Metadata) error {
	_, err := crudobj.CreateAutoScaler(object.AutoScaler{
		TypeMeta: object.TypeMeta{
			Kind:       "AutoScaler",
			APIVersion: "v1",
		},
		ObjectMeta: object.ObjectMeta{
			Name: "ServerlessGatewayAutoScaler",
		},
		Spec: object.AutoScalerSpec{
			Workload: "Pod",
			Template: object.PodTemplate{
				ObjectMeta: object.ObjectMeta{
					Name: "ServerlessGatewayPod",
					Labels: map[string]string{
						options.Usage: options.UsageLabel,
					},
				},
				Spec: object.PodSpec{
					Containers: []object.Container{
						{
							Name:    "ServerlessGatewayContainer",
							Image:   options.GatewayImage,
							Command: []string{meta.Node.Status.Addresses.InternalIP},
						},
					},
				},
			},
			MinReplicas: 1,
			MaxReplicas: 10,
			TargetUtilization: object.UtilizationLimit{
				CPU: &object.CpuUtilizationLimit{
					MinPercentage: 40,
					MaxPercentage: 80,
				},
			},
		},
		Status: nil,
	})

	if err != nil {
		log.Println("[Error]: create auto scaler failed, err:", err.Error())
		return err
	}

	_, err = crudobj.CreateService(object.Service{
		TypeMeta: object.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: object.ObjectMeta{
			Name: "ServerlessGatewayService",
		},
		Spec: object.ServiceSpec{
			Selector: map[string]string{
				options.Usage: options.UsageLabel,
			},
			Ports: []object.ServicePort{
				{
					Protocol:   "TCP",
					Port:       6810,
					TargetPort: 6810,
				},
			},
		},
		Status: nil,
	})

	return nil
}
