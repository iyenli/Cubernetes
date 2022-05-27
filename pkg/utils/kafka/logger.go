package kafka

import (
	"fmt"
	"time"
)

const invalidString = "at offset 0"

func logf(msg string, a ...interface{}) {
	msg = fmt.Sprintf("%v [INFO]: %v", time.Now().Format("2006/01/02 15:04:05"), msg)
	msg = fmt.Sprintf(msg, a...)
	if msg[len(msg)-len(invalidString):] == invalidString {
		return
	}
	fmt.Println(msg)
}
