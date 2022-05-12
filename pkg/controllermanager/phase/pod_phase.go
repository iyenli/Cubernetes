package phase

import (
	"Cubernetes/pkg/object"
)

func NotHandle(p object.PodPhase) bool {
	return p == object.PodCreated || p == object.PodBound
}

func Running(p object.PodPhase) bool {
	return p == object.PodAccepted || p == object.PodPending || p == object.PodRunning
}

func Bad(p object.PodPhase) bool {
	return p == object.PodSucceeded || p == object.PodFailed || p == object.PodUnknown
}
