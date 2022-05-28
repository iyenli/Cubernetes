package phase

import (
	"Cubernetes/pkg/object"
)

func NotHandle(p object.ActorPhase) bool {
	return p == object.ActorCreated || p == object.ActorBound
}

func Running(p object.ActorPhase) bool {
	return p == object.ActorRunning
}

func Failed(p object.ActorPhase) bool {
	return p == object.ActorFailed || p == object.ActorUnknown
}
