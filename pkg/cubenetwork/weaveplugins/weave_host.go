package weaveplugins

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/option"
	"log"
	"net"
	osexec "os/exec"
)

type Host struct {
	IP net.IP
}

// ExposeHost Execute in host, and return an IP in containers' net segment
// you can visit host's port in any container:)
func ExposeHost() (net.IP, error) {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return nil, err
	}

	cmd := osexec.Command(path, "expose")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(output))
	if ip == nil {
		return nil, err
	}

	return ip, nil
}

func CheckSuperUser() error {
	cmd := osexec.Command(option.Sudo, "-s")
	err := cmd.Run()
	if err != nil {
		log.Panicf("Need Sudo: %s\n", err)
		return err
	}

	return nil
}

func InitWeave() error {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave not found.")
		err = InstallWeave()
		if err != nil {
			log.Println("Weave Install failed in adding node to cluster.")
			return err
		}
		if path, err = osexec.LookPath(option.WeaveName); err != nil {
			log.Println("Weave Install but still can't find weave.")
			return err
		}

	}

	err = CheckSuperUser()
	if err != nil {
		return err
	}

	log.Println("Init weave node...")
	// stop weave if weave is running
	cmd := osexec.Command(path, option.Stop)
	err = cmd.Run() // Could failed here

	cmd = osexec.Command(path, option.Launch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Weave add node error: %s, %s\n", err, string(output))
		return err
	}

	return nil
}

// AddNode Called by new node cubelet, it should know its ip and api server's ip.
func AddNode(newHost Host, apiServerHost Host) error {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave not found.")
		err = InstallWeave()
		if err != nil {
			log.Println("Weave Install failed in adding node to cluster.")
			return err
		}
		if path, err = osexec.LookPath(option.WeaveName); err != nil {
			log.Println("Weave Install but still can't find weave.")
			return err
		}

	}

	err = CheckSuperUser()
	if err != nil {
		return err
	}

	// stop weave if weave is running
	cmd := osexec.Command(path, option.Stop)
	err = cmd.Run() // Could failed here

	log.Println("Connecting to peers...")
	cmd = osexec.Command(path, option.Launch, newHost.IP.String(), apiServerHost.IP.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Weave add node error: %s, %v\n", err, string(output))
		return err
	}

	output, err = CheckPeers()
	log.Println("Peers: ", string(output))
	return nil
}

func CheckPeers() ([]byte, error) {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return nil, err
	}

	cmd := osexec.Command(path, "status", "peers")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Weave status error: %s\n", err)
		return nil, err
	}
	return output, nil
}

func CloseNetwork() error {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	cmd := osexec.Command(path, "stop")
	err = cmd.Run()
	if err != nil {
		log.Panicf("Weave reset error: %s\n", err)
		return err
	}

	return nil
}
