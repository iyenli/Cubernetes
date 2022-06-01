package utils

import (
	"Cubernetes/cmd/cuberoot/options"
	"log"
	"os"
	"os/exec"
)

// StartDaemonProcess arg[0]: log, arg[1...n]: args
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

func PreStartMaster() error {
	err := StartDaemonProcess(options.ETCDLOG, options.ETCD)
	if err != nil {
		log.Println("[FATAL] fail to start etcd")
		return err
	}
	err = StartDaemonProcess(options.APISERVERLOG, options.APISERVER)
	if err != nil {
		log.Println("[FATAL] fail to start apiserver")
		return err
	}
	return nil
}

func StartMaster(IP string, NodeUID string) error {
	err := StartDaemonProcess(options.CUBEPROXYLOG, options.CUBEPROXY, IP)
	if err != nil {
		log.Println("[FATAL] fail to start cubeproxy")
		return err
	}
	err = StartDaemonProcess(options.CUBELETLOG, options.CUBELET, NodeUID, IP)
	if err != nil {
		log.Println("[FATAL] fail to start cubelet")
		return err
	}
	err = StartDaemonProcess(options.MANAGERLOG, options.MANAGER, IP)
	if err != nil {
		log.Println("[FATAL] fail to start manager")
		return err
	}
	err = StartDaemonProcess(options.SCHEDULERLOG, options.SCHEDULER, IP)
	if err != nil {
		log.Println("[FATAL] fail to start scheduler")
		return err
	}
	//err = StartDaemonProcess(options.GATEWAYLOG, options.GATEWAY, IP)
	//if err != nil {
	//	log.Println("[FATAL] fail to start gateway")
	//	return err
	//}
	err = StartDaemonProcess(options.BRAINLOG, options.BRAIN, IP)
	if err != nil {
		log.Println("[FATAL] fail to start action brain")
		return err
	}
	return nil
}

func StartSlave(IP, masterIP, NodeUID string) error {
	err := StartDaemonProcess(options.CUBELETLOG, options.CUBELET, NodeUID, IP, masterIP)
	if err != nil {
		log.Println("[FATAL] fail to start cubelet")
		return err
	}
	err = StartDaemonProcess(options.CUBEPROXYLOG, options.CUBEPROXY, IP, masterIP)
	if err != nil {
		log.Println("[FATAL] fail to start cubeproxy")
		return err
	}
	return nil
}
