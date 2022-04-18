/*
	create etcd client and destroy
*/

package etcd_helper

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type ETCDContext struct {
	client *clientv3.Client
}

const etcdTimeout = 3 * time.Second
const etcdAddr = "127.0.0.1:2379"

func newETCDClient() *clientv3.Client {
	log.Println("New etcd client")
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

func closeETCDClient(toClose *clientv3.Client) {
	log.Println("Close etcd client")

	defer func(toClose *clientv3.Client) {
		err := toClose.Close()
		if err != nil {
			log.Panicln("Error: close etcd client failed")
		}
	}(toClose)
}

//func newETCDHealthCheck() (func() error, error) {
//	go wait.PollUntil(time.Second, func() (bool, error) {
//		client, err := newETCD3Client(c.Transport)
//		if err != nil {
//			clientErrMsg.Store(err.Error())
//			return false, nil
//		}
//		clientValue.Store(client)
//		clientErrMsg.Store("")
//		return true, nil
//	}, wait.NeverStop)
//}
