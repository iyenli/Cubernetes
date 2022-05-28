package network

import (
	"Cubernetes/pkg/apiserver/heartbeat"
	network "Cubernetes/pkg/cubelet/network/options"
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"Cubernetes/pkg/utils/localstorage"
	"log"
	"net"
	"os"
)

func InitNodeNetwork(args []string) net.IP {
	log.Println("[INFO]: Starting node network...")
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

	err = weaveplugins.SetWeaveEnv()
	if err != nil {
		log.Panicf("[Error]: Expose host failed, err: %v", err.Error())
		return nil
	}
	ip, err := weaveplugins.ExposeHost()
	if err != nil {
		log.Panicf("[Error]: Expose host failed, err: %v", err.Error())
		return nil
	}
	InitDNSServer()

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

func InitDNSServer() {
	err := os.MkdirAll(network.BaseAddr, 0666)
	if err != nil {
		log.Panicf("[Error]: create dns config dir failed, err: %v", err.Error())
		return
	}

	WriteConfigFile(network.HeadFile, network.HeadFileContent)
}
