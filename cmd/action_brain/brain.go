package main

import brain "Cubernetes/pkg/actionbrain"

func main() {
	brain, err := brain.NewActionBrain()
	if err != nil {
		panic(err)
	}

	brain.Run()
}