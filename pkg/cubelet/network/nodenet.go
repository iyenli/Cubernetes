package network

import (
	"Cubernetes/pkg/apiserver/heartbeat"
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"Cubernetes/pkg/utils/localstorage"
	"log"
	"net"
)

func InitNodeNetwork(args []string) {
	var err error
	if len(args) == 3 {
		// master
		err = weaveplugins.InitWeave()
	} else if len(args) == 4 {
		// slave
		err = weaveplugins.AddNode(weaveplugins.Host{IP: net.ParseIP(args[2])}, weaveplugins.Host{IP: net.ParseIP(args[3])})
	} else {
		panic("Error: too much or little args when start cubelet;")
	}

	if err != nil {
		log.Panicf("Init weave network failed, err: %v", err.Error())
		return
	}
}

func InitNodeHeartbeat() {
	meta, err := localstorage.LoadMeta()
	if err != nil {
		log.Fatal("[FATAL] fail to load node metadata, err: ", err)
		return
	}

	nodenetwork.SetMasterIP(meta.MasterIP)
	heartbeat.InitNode(meta.Node)
}
