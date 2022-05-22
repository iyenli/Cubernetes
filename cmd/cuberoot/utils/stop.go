package utils

import (
	"log"
	"os"
	osexec "os/exec"
	"strconv"
	"strings"
)

func KillDaemonProcess(name string) error {
	Cmd := osexec.Command("pidof", name)
	byteOutput, err := Cmd.CombinedOutput()
	if err != nil {
		log.Printf("[Warn]: No such process named %v\n", name)
		return err
	}

	output := strings.Replace(string(byteOutput), "\n", " ", -1)
	output = strings.Trim(output, " ")
	pids := strings.Split(output, " ")
	log.Printf("[INFO]: Kill %v process soon, name is %v\n", len(pids), name)

	for _, pid := range pids {
		i, err := strconv.Atoi(pid)
		if err != nil {
			log.Println("[Warn]: Invalid PID to parse as int")
			return err
		}
		proc, err := os.FindProcess(i)
		err = proc.Kill()
		if err != nil {
			log.Println("[Warn]: Kill process failed")
			return err
		}
	}

	return nil
}

// ps -ef|grep 进程名关键字|gawk '$0 !~/grep/ {print $2}' |tr -s '\n' ' '
