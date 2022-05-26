package phase

import (
	"Cubernetes/pkg/object"
)

func ActorNotHandle(p object.ActorPhase) bool {
	return p == object.ActorCreated || p == object.ActorBound
}

func ActorRunning(p object.ActorPhase) bool {
	return p == object.ActorRunning
}

func ActorFailed(p object.ActorPhase) bool {
	return p == object.ActorFailed || p == object.ActorUnknown
}
