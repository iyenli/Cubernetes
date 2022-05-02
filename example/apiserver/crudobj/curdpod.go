package main

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"fmt"
)

func main() {
	var pod object.Pod
	pod.APIVersion = "1"
	pod.Kind = "pod"
	pod.Name = "hello4"

	pod, err := crudobj.CreatePod(pod)
	if err != nil {
		return
	}
	fmt.Println("UID:", pod.UID)

	pods, err := crudobj.GetPods()
	if err != nil {
		return
	}
	fmt.Println("Pods:", pods)

	pod.APIVersion = "2"
	pod, err = crudobj.UpdatePod(pod)
	if err != nil {
		return
	}
	fmt.Println("Pod:", pod)

	err = crudobj.DeletePod(pod.UID)
	if err != nil {
		return
	}
}
