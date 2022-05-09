package main

import "Cubernetes/pkg/scheduler"

func main() {
	scheduler := scheduler.NewScheduler()
	scheduler.Run()
}
