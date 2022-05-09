package main

import (
	"Cubernetes/pkg/apiserver/heartbeat"
	"Cubernetes/pkg/object"
	"fmt"
	"time"
)

func main() {
	node := object.Node{}
	node.APIVersion = "1"
	node.Kind = "node"
	node.Name = "hello6"
	node.UID = "5533241f-a05f-4808-97dc-08f2248fccfb"
	node.Status = &object.NodeStatus{}
	node.Status.Condition.Ready = true
	fmt.Println(node)

	heartbeat.InitNode(node)
	time.Sleep(100 * time.Second)
}
