/*
	create etcd Client and destroy
*/

package etcd_helper

import (
	cubeconfig "Cubernetes/config"
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type ETCDContext struct {
	Client *clientv3.Client
}

func NewETCDClient() *clientv3.Client {
	log.Println("New etcd Client")
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{cubeconfig.ETCDAddr},
		DialTimeout: cubeconfig.ETCDTimeout,
	})

	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return nil
	}

	return client
}

func CloseETCDClient(toClose *clientv3.Client) {
	log.Println("Close etcd Client")

	defer func(toClose *clientv3.Client) {
		err := toClose.Close()
		if err != nil {
			log.Panicln("Error: close etcd Client failed")
		}
	}(toClose)
}

func ETCDHealthCheck(ctx *ETCDContext) bool {
	ticker := time.NewTicker(cubeconfig.ETCDTimeout)
	health := make(chan bool)

	go func() {
		_, err := NewETCDClient().KV.Get(context.TODO(), "HealthCheck")
		if err != nil {
			health <- false
			return
		}
		health <- true
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
