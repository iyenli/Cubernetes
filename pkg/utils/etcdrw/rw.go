package etcdrw

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func GetObj(path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	res, err := client.Get(ctx, path)
	cancel()
	if err != nil {
		log.Panicf("fail to get object from etcd, path: %v, err: %v\n", path, err)
		return nil, err
	}
	if res.Count == 0 {
		log.Printf("no objects found in etcd, path: %v\n", path)
		return nil, nil
	}
	return res.Kvs[0].Value, nil
}

func GetObjs(prefix string) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	res, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Panicf("fail to get objects from etcd, prefix: %v, err: %v\n", prefix, err)
		return nil, err
	}
	if res.Count == 0 {
		log.Printf("no objects found in etcd, prefix: %v\n", prefix)
		return nil, nil
	}
	var ret [][]byte
	for _, kv := range res.Kvs {
		ret = append(ret, kv.Value)
	}
	return ret, nil
}

func PutObj(path string, obj string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	_, err := client.Put(ctx, path, obj)
	cancel()
	if err != nil {
		log.Panicf("fail to put object into etcd, path: %v, err: %v\n", path, err)
	}
	return err
}

func DelObj(path string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	_, err := client.Delete(ctx, path)
	cancel()
	if err != nil {
		log.Panicf("fail to delete object from etcd, path: %v, err: %v\n", path, err)
	}
	return err
}
