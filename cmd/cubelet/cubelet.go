package main

import (
	"Cubernetes/pkg/cubelet"
)

func main() {
	cubeletInstance := cubelet.NewCubelet()
	cubeletInstance.Run()
}
