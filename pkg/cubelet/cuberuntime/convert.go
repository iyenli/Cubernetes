package cuberuntime

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	"sort"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func toRuntimeAPIImageSpec(imageSpec cubecontainer.ImageSpec) *runtimeapi.ImageSpec {
	var annotations = make(map[string]string)
	if imageSpec.Annotations != nil {
		for _, a := range imageSpec.Annotations {
			annotations[a.Name] = a.Value
		}
	}
	return &runtimeapi.ImageSpec{
		Image:       imageSpec.Image,
		Annotations: annotations,
	}
}

func toCubeContainerImageSpec(image *runtimeapi.Image) cubecontainer.ImageSpec {
	var annotations []cubecontainer.Annotation

	if image.Spec != nil && len(image.Spec.Annotations) > 0 {
		annotationKeys := make([]string, 0, len(image.Spec.Annotations))
		for k := range image.Spec.Annotations {
			annotationKeys = append(annotationKeys, k)
		}
		sort.Strings(annotationKeys)
		for _, k := range annotationKeys {
			annotations = append(annotations, cubecontainer.Annotation{
				Name:  k,
				Value: image.Spec.Annotations[k],
			})
		}
	}

	return cubecontainer.ImageSpec{
		Image:       image.Id,
		Annotations: annotations,
	}
}
