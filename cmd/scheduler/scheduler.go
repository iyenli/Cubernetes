package main

import "Cubernetes/pkg/scheduler"

func main() {
	newScheduler := scheduler.NewScheduler()
	newScheduler.Run()
}
