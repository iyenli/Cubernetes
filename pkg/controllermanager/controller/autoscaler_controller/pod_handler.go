package autoscaler_controller

import "Cubernetes/pkg/object"

func (asc *autoScalerController) handlePodCreate(pod *object.Pod) error {
	uid, ok := pod.Labels[lowerReplocaSetParentUIDLabel]
	if ok {
		as, found := asc.asInformer.GetAutoScaler(uid)
		if found {
			return asc.checkAndUpdateAutoScalerStatus(as)
		}
	}
	return nil
}