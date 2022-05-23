package kafka

import "fmt"

func logf(msg string, a ...interface{}) {
	msg = fmt.Sprintf("[INFO]: %v\n", msg)
	fmt.Printf(msg, a...)
}
