package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetPod(UID string) (object.Pod, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/pod/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.Pod{}, err
	}

	var pod object.Pod
	err = json.Unmarshal(body, &pod)
	if err != nil {
		log.Println("fail to parse Pod")
		return object.Pod{}, err
	}

	return pod, nil
}

func GetPods() ([]object.Pod, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/pods"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var pods []object.Pod
	err = json.Unmarshal(body, &pods)
	if err != nil {
		log.Println("fail to parse Pods")
		return nil, err
	}

	return pods, nil
}

func SelectPods(selectors map[string]string) ([]object.Pod, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/pods"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var pods []object.Pod
	err = json.Unmarshal(body, &pods)
	if err != nil {
		log.Println("fail to parse Pods")
		return nil, err
	}

	return pods, nil
}

func CreatePod(pod object.Pod) (object.Pod, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/pod"

	body, err := postRequest(url, pod)
	if err != nil {
		log.Println("postRequest fail")
		return pod, err
	}

	var newPod object.Pod
	err = json.Unmarshal(body, &newPod)
	if err != nil {
		log.Println("fail to parse Pod")
		return pod, err
	}

	return newPod, nil
}

func UpdatePod(pod object.Pod) (object.Pod, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/pod/" + pod.UID

	body, err := putRequest(url, pod)
	if err != nil {
		log.Println("putRequest fail")
		return pod, err
	}

	var newPod object.Pod
	err = json.Unmarshal(body, &newPod)
	if err != nil {
		log.Println("fail to parse Pod")
		return pod, err
	}

	return newPod, nil
}

func UpdatePodStatus(UID string, status object.PodStatus) (object.Pod, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/pod/status/" + UID

	pod := object.Pod{}
	pod.UID = UID
	pod.Status = &status

	body, err := putRequest(url, pod)
	if err != nil {
		log.Println("putRequest fail")
		return pod, err
	}

	var newPod object.Pod
	err = json.Unmarshal(body, &newPod)
	if err != nil {
		log.Println("fail to parse Pod")
		return pod, err
	}

	return newPod, nil
}

func DeletePod(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/pod/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("postRequest fail")
		return err
	}

	return nil
}
