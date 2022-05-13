package utils

import (
	"log"
	osexec "os/exec"
)

const ETCD = "etcdctl"

func ClearData() error {
	path, err := osexec.LookPath(ETCD)
	if err != nil {
		log.Panicf("Command not found: %v", err.Error())
		return err
	}

	cmd := osexec.Command(path, "del", "/", "--prefix")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Clear error: output: %v, err: %v", string(output), err.Error())
		return err
	}

	log.Println("Clear ETCD KVs number:", string(output))
	return nil
}
