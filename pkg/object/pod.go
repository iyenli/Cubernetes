package object

type Pod struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
		UID       string            `json:"uid"`
		Labels    map[string]string `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		Containers []struct {
			Name  string `json:"name"`
			Image string `json:"image"`
			Ports []struct {
				ContainerPort int    `json:"containerPort"`
				Name          string `json:"name"`
				Protocol      string `json:"protocol"`
			} `json:"ports"`
			Env          []string `json:"env"`
			VolumeMounts []struct {
				Name      string `json:"name"`
				MountPath string `json:"mountPath"`
			} `json:"volumeMounts"`
			Args []string
		} `json:"containers"`
		Volumes []struct {
			Name string `json:"name"`
		} `json:"volumes"`
	} `json:"spec"`
	Status string `json:"status"`
}
