package object

func ComputePodSpecChange(new *Pod, old *Pod) bool {
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

	// Spec change
	if len(new.Spec.Containers) != len(old.Spec.Containers) {
		return true
	}
	for i, oldC := range old.Spec.Containers {
		if ComputeContainerSpecChange(&new.Spec.Containers[i], &oldC) {
			return true
		}
	}

	if len(new.Spec.Volumes) != len(old.Spec.Volumes) {
		return true
	}
	for i, oldV := range old.Spec.Volumes {
		newV := new.Spec.Volumes[i]
		if newV.HostPath != oldV.HostPath || newV.Name != oldV.Name {
			return true
		}
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
