package main

import (
	"Cubernetes/cmd/apiserver/heartbeat"
	"Cubernetes/cmd/apiserver/httpserver"
	"Cubernetes/cmd/apiserver/httpserver/restful"
	"Cubernetes/pkg/cubenetwork/servicenetwork"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"log"
	"sync"
	"time"
)

func main() {
	etcdrw.Init()
	defer etcdrw.Free()

	time.Sleep(time.Second)
	updateNodeReadyState()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		heartbeat.ListenHeartbeat()
	}()
	go func() {
		defer wg.Done()
		httpserver.Run()
	}()

	time.Sleep(time.Second)
	restful.ClusterIPAllocator = servicenetwork.NewClusterIPAllocator()

	log.Println("[INFO]: Cluster IP Allocator init, api server running...")
	wg.Wait()
}

func updateNodeReadyState() {
	nodes, err := etcdrw.GetObjs(object.NodeEtcdPrefix)
	if err != nil {
		log.Fatal("[FATAL] fail to get Nodes from etcd")
	}

	for _, buf := range nodes {
		var node object.Node
		err = json.Unmarshal(buf, &node)
		if err != nil {
			log.Fatal("[FATAL] fail to parse Node")
		}

		if node.Status == nil || node.Status.Condition.Ready == false {
			continue
		}

		node.Status.Condition.Ready = false
		newBuf, err := json.Marshal(node)
		if err != nil {
			log.Fatal("[FATAL] marshal to parse Node")
		}

		err = etcdrw.PutObj(object.NodeEtcdPrefix+node.UID, string(newBuf))
		if err != nil {
			log.Fatal("[FATAL] fail to put Node into etcd")
		}
	}
}
