package main

import (
	"Cubernetes/pkg/cubelet"
)

func main() {
	cubelet := cubelet.NewCubelet()
	cubelet.Run()
}
