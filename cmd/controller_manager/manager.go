package main

import "Cubernetes/pkg/controllermanager"

func main() {
	cm := controllermanager.NewControllerManager()
	cm.Run()
}
