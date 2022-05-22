package network

import (
	network "Cubernetes/pkg/cubelet/network/options"
	"log"
	"os"
	"os/exec"
)

func WriteConfigFile(filename, content string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Panicf("[Error]: create dns config file failed, err: %v", err.Error())
		return
	}

	_, err = f.Write([]byte(content))
	if err != nil {
		log.Panicf("[Error]: write dns config file failed, err: %v", err.Error())
		return
	}

	cmd := exec.Command(network.Resolve, "-u")
	err = cmd.Run()
	if err != nil {
		log.Panicf("[Error]: reset resolve config failed, err: %v", err.Error())
		return
	}
}
