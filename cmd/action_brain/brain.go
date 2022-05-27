package main

import brain "Cubernetes/pkg/actionbrain"

func main() {
	brainRuntime, err := brain.NewActionBrain()
	if err != nil {
		panic(err)
	}

	brainRuntime.Run()
}
