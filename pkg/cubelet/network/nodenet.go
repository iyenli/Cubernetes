package network

import (
	"Cubernetes/pkg/apiserver/heartbeat"
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"Cubernetes/pkg/utils/localstorage"
	"log"
	"net"
	"time"
)

func InitNodeNetwork(args []string) net.IP {
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
		return nil
	}

	// wait for weave starting
	time.Sleep(12 * time.Second)
	ip, err := weaveplugins.ExposeHost()
	if err != nil {
		log.Panicf("Expose host failed, err: %v", err.Error())
		return nil
	}

	log.Println("[INFO] weave IP Allocated:", ip.String())
	return ip
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
