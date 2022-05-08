package main

import (
	"os"
	"os/signal"
	"time"
)

// deprecated
func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)

	go func() {
		defer func() {
			println("Oops!")
		}()
		for true {
			println("Hi guys")
			time.Sleep(time.Second)
		}
	}()

	<-c
	println("Ended.")
}

//var c chan os.Signal

//func main() {
//	c := make(chan os.Signal)
//	signal.Notify(c, os.Kill, os.Interrupt)
//
//	go func() {
//		<-c
//		println("fake defer")
//		os.Exit(1)
//	}()
//
//	for true {
//		println("working...")
//		time.Sleep(time.Second)
//	}
//}
