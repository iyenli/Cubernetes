package utils

import (
	"log"
	"os"
	"os/exec"
)

// StartDaemonProcess arg[0]: log, arg[1...n]args
func StartDaemonProcess(args ...string) error {
	_, err := os.Stat(args[1])
	if err != nil {
		log.Panicf("Command not found: %v", err.Error())
		return err
	}

	server := exec.Command(args[1], args[2:]...)
	// Don't close here:)
	stdout, err := os.OpenFile(args[0], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		log.Println(os.Getpid(), ": open log file error", err)
	}
	server.Stderr = stdout
	server.Stdout = stdout
	err = server.Start()

	if err != nil {
		return err
	}
	return nil
}
