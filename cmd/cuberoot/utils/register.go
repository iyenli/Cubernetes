package utils

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/localstorage"
	"log"
	"os"
	"time"
)

func autoFillNode(node *object.Node) {
	node.Spec.Info.CubeVersion = cubeconfig.CubeVersion
	node.Spec.Info.DeviceName, _ = os.Hostname()
	node.Status.Condition.Ready = false
	node.Status.Condition.DiskPressure = false
	node.Status.Condition.MemoryPressure = false
	node.Status.Condition.OutOfDisk = false
}

func RegisterAsMaster(node object.Node) error {
	node.Spec.Type = object.Master
	node.Kind = cubeconfig.KindNode
	autoFillNode(&node)

	node, err := crudobj.CreateNode(node)
	if err != nil {
		log.Println("[FATAL] fail to create node, err: ", err)
		return err
	}

	err = localstorage.SaveMeta(localstorage.Metadata{
		Node:     node,
		MasterIP: node.Status.Addresses.InternalIP,
	})
	if err != nil {
		log.Println("[FATAL] fail to save node metadata, err: ", err)
	}
	return err
}

func RegisterAsSlave(node object.Node, masterIP string) error {
	node.Spec.Type = object.Slave
	node.Kind = cubeconfig.KindNode
	autoFillNode(&node)

	node, err := crudobj.CreateNode(node)
	if err != nil {
		log.Println("[FATAL] fail to create node, err: ", err)
		return err
	}

	err = localstorage.SaveMeta(localstorage.Metadata{
		Node:     node,
		MasterIP: masterIP,
	})
	if err != nil {
		log.Println("[FATAL] fail to save node metadata, err: ", err)
	}
	return err
}

func StartFromRegistry() {
	meta, err := localstorage.LoadMeta()
	if err != nil {
		log.Fatal("[FATAL] fail to load node metadata, err: ", err)
	}

	if meta.Node.Spec.Type == object.Master {
		log.Println("Starting as master, this may take 10s")

		err = PreStartMaster()
		if err != nil {
			log.Fatal("[FATAL] fail to prestart master processes, err: ", err)
		}

		time.Sleep(10 * time.Second)

		err = StartMaster(meta.Node.Status.Addresses.InternalIP, meta.Node.UID)
		if err != nil {
			log.Fatal("[FATAL] fail to start master processes, err: ", err)
		}

		log.Printf("Master node launched successfully\n"+
			"To join Cubernetes cluster, execute:\n"+
			"\tcuberoot join %s -f [node config file]\n", meta.Node.Status.Addresses.InternalIP)

	} else {
		err = StartSlave(meta.Node.Status.Addresses.InternalIP, meta.MasterIP, meta.Node.UID)
		if err != nil {
			log.Fatal("[FATAL] fail to start slave processes, err: ", err)
		}
	}
}
