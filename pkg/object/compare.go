package object

import "log"

func ComputeObjectMetaChange(new *ObjectMeta, old *ObjectMeta) bool {
	if new.UID != old.UID {
		panic("Compute 2 pod change with different uid")
	}

	// name change?
	if new.Name != old.Name {
		return true
	}

	// label change?
	if len(new.Labels) != len(old.Labels) {
		return true
	}
	for k, oldV := range old.Labels {
		newV, ok := new.Labels[k]
		if !ok || oldV != newV {
			return true
		}
	}

	return false
}

func ComputePodSpecChange(new *PodSpec, old *PodSpec) bool {

	// Spec change
	if len(new.Containers) != len(old.Containers) {
		return true
	}
	for i, oldC := range old.Containers {
		if ComputeContainerSpecChange(&new.Containers[i], &oldC) {
			return true
		}
	}

	if len(new.Volumes) != len(old.Volumes) {
		return true
	}
	for i, oldV := range old.Volumes {
		newV := new.Volumes[i]
		if newV.HostPath != oldV.HostPath || newV.Name != oldV.Name {
			return true
		}
	}

	return false
}

func ComputeReplicaSetSpecChange(new *ReplicaSetSpec, old *ReplicaSetSpec) bool {
	if new.Replicas != old.Replicas {
		return true
	}

	if len(new.Selector) != len(old.Selector) {
		return true
	}
	for k, oldV := range old.Selector {
		newV, ok := new.Selector[k]
		if !ok || newV != oldV {
			return true
		}
	}

	if ComputeObjectMetaChange(&new.Template.ObjectMeta, &old.Template.ObjectMeta) {
		return true
	}

	if ComputePodSpecChange(&new.Template.Spec, &old.Template.Spec) {
		return true
	}

	return false
}

func ComputeContainerSpecChange(new *Container, old *Container) bool {
	// check basic info
	if new.Name != old.Name || new.Image != old.Image {
		return true
	}

	// check Command & Args
	if len(new.Command) != len(old.Command) {
		return true
	}
	for i, c := range old.Command {
		if new.Command[i] != c {
			return true
		}
	}

	if len(new.Args) != len(old.Args) {
		return true
	}
	for i, a := range old.Args {
		if new.Args[i] != a {
			return true
		}
	}

	// check Resource limits
	if new.Resources != nil && old.Resources != nil {
		if new.Resources.Cpus != old.Resources.Cpus ||
			new.Resources.Memory != old.Resources.Memory {
			return true
		}
	} else if new.Resources != nil || old.Resources != nil {
		return true
	}

	// check VolumeMounts
	if len(new.VolumeMounts) != len(old.VolumeMounts) {
		return true
	}
	for i, oldM := range old.VolumeMounts {
		newM := new.VolumeMounts[i]
		if newM.Name != oldM.Name || newM.MountPath != oldM.MountPath {
			return true
		}
	}

	// check Ports
	if len(new.Ports) != len(old.Ports) {
		return true
	}
	for i, oldP := range old.Ports {
		newP := new.Ports[i]
		if newP.Name != oldP.Name || newP.HostPort != oldP.HostPort || newP.HostIP != oldP.HostIP ||
			newP.ContainerPort != oldP.ContainerPort || newP.Protocol != oldP.Protocol {
			return true
		}
	}

	return false
}

func MatchLabelSelector(selector map[string]string, labels map[string]string) bool {
	for k, v := range selector {
		podV, ok := labels[k]
		if !ok || podV != v {
			return false
		}
	}
	return true
}

// ComputePodNetworkChange Just check label and ip
func ComputePodNetworkChange(new *Pod, old *Pod) bool {
	for k, oldV := range old.ObjectMeta.Labels {
		newV, ok := new.ObjectMeta.Labels[k]
		if !ok || newV != oldV {
			return true
		}
	}

	if !new.Status.IP.Equal(old.Status.IP) {
		return true
	}

	if old.Status.Phase == PodRunning && new.Status.Phase != PodRunning ||
		new.Status.Phase == PodRunning && old.Status.Phase != PodRunning {
		log.Println("Some pod crashed or restarted, reset service")
		return true
	}

	return false
}

// ComputeServiceCriticalChange if true, we have to reset iptables related to this service
func ComputeServiceCriticalChange(new *Service, old *Service) bool {
	// Selector: affect pods selected
	if len(new.Spec.Selector) != len(old.Spec.Selector) {
		return true
	}
	for k, oldV := range old.Spec.Selector {
		newV, ok := new.Spec.Selector[k]
		if !ok || oldV != newV {
			return true
		}
	}

	// port affect iptables directly
	if len(new.Spec.Ports) != len(old.Spec.Ports) {
		return true
	}
	for _, oldPort := range old.Spec.Ports {
		isSame := false
		for _, newPort := range new.Spec.Ports {
			if oldPort.Protocol == newPort.Protocol &&
				oldPort.Port == newPort.Port &&
				oldPort.TargetPort == newPort.TargetPort {
				isSame = true
				break
			}
		}
		if !isSame {
			return true
		}
	}

	// cluster ip affect iptables directly
	if new.Spec.ClusterIP != old.Spec.ClusterIP {
		return true
	}

	// We don't care endpoints, because every proxy judges pods independently
	// TODO: is ingress critical?
	return false
}

func ComputeDNSCriticalChange(new *Dns, old *Dns) bool {
	if new.Spec.Host != old.Spec.Host {
		return true
	}
	if len(new.Spec.Paths) != len(old.Spec.Paths) {
		return true
	}

	for key, val := range new.Spec.Paths {
		if tmp, exist := old.Spec.Paths[key]; exist {
			if tmp.ServicePort != val.ServicePort ||
				tmp.ServiceUID != val.ServiceUID {
				return true
			}
		} else {
			return true
		}
	}

	return false
}
