package objconfig

import "gopkg.in/yaml.v3"

type Config interface {
	GetKind() string
}

type PodConfig struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string            `yaml:"name"`
		Namespace string            `yaml:"namespace"`
		UID       string            `yaml:"uid"`
		Labels    map[string]string `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Containers []struct {
			Name  string `yaml:"name"`
			Image string `yaml:"image"`
			Ports []struct {
				ContainerPort int    `yaml:"containerPort"`
				Name          string `yaml:"name"`
				Protocol      string `yaml:"protocol"`
			} `yaml:"ports"`
			Env          []string `yaml:"env"`
			VolumeMounts []struct {
				Name      string `yaml:"name"`
				MountPath string `yaml:"mountPath"`
			} `yaml:"volumeMounts"`
			Args []string
		} `yaml:"containers"`
		Volumes []struct {
			Name string `yaml:"name"`
		} `yaml:"volumes"`
	} `yaml:"spec"`
}

func (pc PodConfig) GetKind() string {
	return pc.Kind
}

func ParseConfig(file []byte) Config {
	configMap := make(map[interface{}]interface{})
	err := yaml.Unmarshal(file, &configMap)
	if err != nil {
		return nil
	}

	var config Config
	switch configMap["kind"] {
	case "Pod":
		var podConfig PodConfig
		err = yaml.Unmarshal(file, &podConfig)
		if err != nil {
			return nil
		}
		config = podConfig
	default:
		return nil
	}
	return config
}
