package weaveplugins

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/option"
	"log"
	"net"
	"os"
	osexec "os/exec"
	"strings"
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
	combinedOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("[INFO]: weave expose failed")
		return nil, err
	}

	output := strings.Trim(string(combinedOutput), "\n")
	ip := net.ParseIP(string(output))
	if ip == nil {
		log.Println("[INFO]: Parse host weave ip failed")
		return nil, err
	}

	log.Println("[INFO]: Allocate Host IP here, IP is", ip.String())
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

	log.Println("[INFO]: Init weave node...")
	// stop weave if weave is running
	cmd := osexec.Command(path, option.Stop)
	err = cmd.Run() // Could fail here
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

	log.Println("Connecting to peers...")
	cmd = osexec.Command(path, option.Launch, newHost.IP.String(), apiServerHost.IP.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Weave add node error: %s, %v\n", err, string(output))
		return err
	}

	output, err = CheckPeers()
	if err != nil {
		return err
	}
	log.Println("[INFO]: Peers: ", string(output))
	return nil
}

func CheckPeers() ([]byte, error) {
	path, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("[Error]: Weave Not found.")
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

// SetWeaveEnv Execute in host, set docker env
// FIXME: When network crash, consider this function:)
func SetWeaveEnv() error {
	_, err := osexec.LookPath(option.WeaveName)
	if err != nil {
		log.Println("[Error]: Weave Not found")
		return err
	}

	err = os.MkdirAll(option.InitScriptFileDir, 0666)
	if err != nil {
		log.Println("[Error]: Mkdir failed")
		return err
	}

	f, err := os.Create(option.InitScriptFile)
	if err != nil {
		log.Println("[Error]: Mkdir failed")
		return err
	}

	_, err = f.Write([]byte(option.InitScript))
	if err != nil {
		log.Println("[Error]: Create script failed")
		return err
	}

	cmd := osexec.Command("/bin/bash", option.InitScriptFile)
	err = cmd.Run()
	if err != nil {
		log.Panicf("Weave init env error: %s\n", err)
		return err
	}

	return nil
}
