package host

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/weave"
	"log"
	"net"
	osexec "os/exec"
)

const (
	weaveName = "weave"
	launch    = "launch"
	sudo      = "sudo"
)

type Host struct {
	IP net.IP
}

func CheckSuperUser() error {
	cmd := osexec.Command(sudo, "-s")
	err := cmd.Run()
	if err != nil {
		log.Panicf("Need Sudo: %s\n", err)
		return err
	}

	return nil
}

func InitWeave(server Host) error {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave not found.")
		err = weave.InstallWeave()
		if err != nil {
			log.Println("Weave Install failed in adding node to cluster.")
			return err
		}
		if path, err = osexec.LookPath(weaveName); err != nil {
			log.Println("Weave Install but still can't find weave.")
			return err
		}

	}

	err = CheckSuperUser()
	if err != nil {
		return err
	}

	cmd := osexec.Command(path, launch, server.IP.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Weave add node error: %s, %s\n", err, string(output))
		return err
	}

	return nil
}

// AddNode Called by new node cubelet, it should know its ip and api server's ip.
func AddNode(newHost Host, apiServerHost Host) error {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave not found.")
		err = weave.InstallWeave()
		if err != nil {
			log.Println("Weave Install failed in adding node to cluster.")
			return err
		}
		if path, err = osexec.LookPath(weaveName); err != nil {
			log.Println("Weave Install but still can't find weave.")
			return err
		}

	}

	err = CheckSuperUser()
	if err != nil {
		return err
	}

	cmd := osexec.Command(path, launch, newHost.IP.String(), apiServerHost.IP.String())
	err = cmd.Run()
	if err != nil {
		log.Panicf("Weave add node error: %s\n", err)
		return err
	}

	return nil
}

func CheckPeers() ([]byte, error) {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return nil, err
	}

	cmd := osexec.Command(path, "status", "peers")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Weave add node error: %s\n", err)
		return nil, err
	}
	return output, nil
}
