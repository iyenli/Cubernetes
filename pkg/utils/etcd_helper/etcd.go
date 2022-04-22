/*
	create etcd Client and destroy
*/

package etcd_helper

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type ETCDContext struct {
	Client *clientv3.Client
}

const etcdTimeout = 3 * time.Second
const etcdAddr = "127.0.0.1:2379"

func NewETCDClient() *clientv3.Client {
	log.Println("New etcd Client")
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: etcdTimeout,
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
	ticker := time.NewTicker(etcdTimeout)
	health := make(chan bool)

	go func() {
		_, err := NewETCDClient().KV.Get(context.TODO(), "HealthCheck")
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
