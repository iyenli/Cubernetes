package etcdrw

import (
	cubeconfig "Cubernetes/config"
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

var client *clientv3.Client

func Init() {
	log.Println("initializing etcd client...")

	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{cubeconfig.ETCDAddr},
		DialTimeout: cubeconfig.ETCDTimeout,
	})

	if err != nil {
		log.Fatalf("fail to initialize etcd client, err: %v\n", err)
		return
	}

	log.Println("etcd client initialized")
}

func Free() {
	log.Println("closing etcd client...")
	err := client.Close()
	if err != nil {
		log.Panicf("fail to close etcd client, err:%v\n", err)
	}
	log.Println("etcd client closed")
}

func CheckHealth() bool {
	ticker := time.NewTicker(cubeconfig.ETCDTimeout)
	health := make(chan bool)

	go func() {
		_, err := client.KV.Get(context.TODO(), "HealthCheck")
		if err != nil {
			health <- false
			return
		}
		health <- true
		return
	}()

	for {
		select {
		case res := <-health:
			return res
		case <-ticker.C:
			log.Println("ETCD Health check time out")
			return false
		}
	}
}
