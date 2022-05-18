package main

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"encoding/json"
	"fmt"
)

func main() {
	var pod object.Pod
	pod.APIVersion = "1"
	pod.Kind = object.KindPod
	pod.Name = "hello4"

	pod.Labels = make(map[string]string)
	pod.Labels["app"] = "hello"
	pod.Labels["container"] = "1"

	selector := pod.Labels

	pod, err := crudobj.CreatePod(pod)
	if err != nil {
		return
	}
	fmt.Println("UID:", pod.UID)

	selectedPods, err := crudobj.SelectPods(selector)
	if err != nil {
		return
	}
	fmt.Println("Selected Pods:", selectedPods)

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

	status := object.PodStatus{
		IP:                  nil,
		Phase:               object.PodCreated,
		ActualResourceUsage: nil,
	}
	pod, err = crudobj.UpdatePodStatus(pod.UID, status)
	buf, err := json.Marshal(pod)
	fmt.Println("Pod:", string(buf))

	err = crudobj.DeletePod(pod.UID)
	if err != nil {
		return
	}
}
