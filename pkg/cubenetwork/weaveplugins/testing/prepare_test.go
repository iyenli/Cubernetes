package testing

import (
	"log"
	osexec "os/exec"
)

const (
	weaveName = "weave"
	launch    = "launch"
	sudo      = "sudo"
)

// delete containers and close weave
func PrepareTest() error {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	cmd := osexec.Command(path, "stop")
	err = cmd.Run()
	if err != nil {
		log.Panicf("Weave stop error: %s\n", err)
		return err
	}

	cmd = osexec.Command("docker", "rm", "$(docker ps -a -q)")
	err = cmd.Run()
	if err != nil {
		log.Panicf("rm docker error: %s\n", err)
		return err
	}

	return nil
}
