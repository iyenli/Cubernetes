package testing

import (
	"log"
	osexec "os/exec"
	"strings"
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

	cmd = osexec.Command("docker", "ps", "-aq")
	byteOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("rm docker error: %s\n", err)
		return err
	}

	// Warning: If you have docker running, be careful of this test
	if len(string(byteOutput)) != 0 {
		output := strings.ReplaceAll(string(byteOutput), "\n", " ")

		cmd = osexec.Command("docker", "stop", output)
		err = cmd.Run()

		cmd = osexec.Command("docker", "rm", output)
		err = cmd.Run()
	}
	return nil
}

func RunContainer() string {
	cmd := osexec.Command("docker", "run", "-d", "-ti", "weaveworks/ubuntu")
	byteOutput, _ := cmd.CombinedOutput()

	return string(byteOutput)
}
