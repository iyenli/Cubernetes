package actor_runtime

import (
	"Cubernetes/pkg/object"
	"strings"
)

const (
	cubePrefix           = "c8s"
	sandboxContainerName = "ACTORSandbox"
	containerName        = "actor"
	nameDelimiter        = "_"
	nameUUIDLen          = 8
)

func makeSandboxName(actor *object.Actor) string {
	return strings.Join([]string{
		cubePrefix,
		actor.Name,
		sandboxContainerName,
		actor.UID[:nameUUIDLen],
	}, nameDelimiter)
}

func makeContainerName(actor *object.Actor) string {
	return strings.Join([]string{
		cubePrefix,
		actor.Name,
		containerName,
		actor.UID[:nameUUIDLen],
	}, nameDelimiter)
}
