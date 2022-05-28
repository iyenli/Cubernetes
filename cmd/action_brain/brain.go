package main

import (
	brain "Cubernetes/pkg/actionbrain"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("[FATAL] Lack arguments")
	}

	brainRuntime, err := brain.NewActionBrain(os.Args[1])
	if err != nil {
		panic(err)
	}

	brainRuntime.Run()
}
