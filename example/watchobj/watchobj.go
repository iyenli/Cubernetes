package main

import (
	"Cubernetes/pkg/watchobj"
	"fmt"
)

// simple example of use
func main() {
	ch, cancel := watchobj.WatchObj("/apis/watch/pod/hello:e0a77a11-f736-4f5f-934e-f1f0a3c39172")
	i := 1
	for str := range ch {
		if i > 5 {
			break
		}
		i++
		fmt.Println(str)
	}
	// cancel watching
	cancel()
}
