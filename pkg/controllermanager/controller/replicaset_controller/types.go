package replicaset_controller

import "Cubernetes/pkg/object"

/*
	** Only Support: **
	cubectl autoscale service service_name --min=2 --max=10 --cpu-percent 50 --memory 4g
	- if no memory or cpu percent, discard this metric
	- default min value is 1
	- No default max value
	- for service only
*/

import (
	"time"
)

/**
Scale: Expected and actual pods
*/

// Scale represents a scaling request for a resource.
type Scale struct {
	object.TypeMeta
	object.ObjectMeta

	// defines the behavior of the scale.
	Spec ScaleSpec

	// current status of the scale.
	Status ScaleStatus
}

// ScaleSpec describes the attributes of a scale subresource
type ScaleSpec struct {
	// desired number of instances for the scaled object.
	Replicas int `json:"replicas,omitempty"`
}

// ScaleStatus represents the current status of a scale subresource.
type ScaleStatus struct {
	// actual number of observed instances of the scaled object.
	Replicas int `json:"replicas"`
}

/*
HorizontalPodAutoscaler Object
*/

// HorizontalPodAutoscaler configuration of a horizontal pod autoscaler.
type HorizontalPodAutoscaler struct {
	object.TypeMeta
	object.ObjectMeta

	// behavior of autoscaler.
	Spec HorizontalPodAutoscalerSpec

	// current information about the autoscaler.
	Status HorizontalPodAutoscalerStatus
}

type Utilization struct {
	// fraction of the requested CPU that should be utilized/used,
	// e.g. 70 means that 70% of the requested CPU should be in use.
	CPUTargetPercentage float64

	MemoryTargetBytes int64
}

// HorizontalPodAutoscalerSpec specification of a horizontal pod autoscaler.
type HorizontalPodAutoscalerSpec struct {
	// reference to Scale subresource; horizontal pod autoscaler will learn the current resource
	// consumption from its status,and will set the desired number of pods by modifying its spec.
	ScaleRef *object.Service

	// lower limit for the number of pods that can be set by the autoscaler, default 1.
	MinReplicas int

	// upper limit for the number of pods that can be set by the autoscaler.
	// It cannot be smaller than MinReplicas.
	MaxReplicas int

	// target average CPU utilization (represented as a percentage of requested CPU) over all the pods;
	// if not specified it defaults to the target CPU utilization at 80% of the requested resources.
	TargetUtilization Utilization
}

// HorizontalPodAutoscalerStatus current status of a horizontal pod autoscaler
type HorizontalPodAutoscalerStatus struct {
	// most recent generation observed by this autoscaler.
	ObservedGeneration int64

	// last time the HorizontalPodAutoscaler scaled the number of pods;
	// used by the autoscaler to control how often the number of pods is changed.
	LastScaleTime time.Time

	// current number of replicas of pods managed by this autoscaler.
	CurrentReplicas int

	// desired number of replicas of pods managed by this autoscaler.
	DesiredReplicas int

	CurrentUtilization *Utilization
}

type HorizontalPodAutoscalerList struct {
	object.TypeMeta
	object.ObjectMeta

	// list of horizontal pod autoscaler objects.
	Items []HorizontalPodAutoscaler
}
