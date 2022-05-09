package utils

import (
	osexec "os/exec"
	"strings"
)

// KillDaemonProcess TODO: A more elegant implement?
func KillDaemonProcess(name string) error {
	Cmd := osexec.Command("pidof", name)
	byteOutput, err := Cmd.CombinedOutput()
	if err != nil {

		return err
	}

	output := strings.Replace(string(byteOutput), "\n", " ", -1)
	output = strings.Trim(output, " ")

	if len(output) != 0 {
		cmd := osexec.Command("kill", "-9", output)
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// ps -ef|grep 进程名关键字|gawk '$0 !~/grep/ {print $2}' |tr -s '\n' ' '
