package etcdrw

import (
	"context"
	"go.etcd.io/etcd/clientv3"
)

var watchMap map[string]context.CancelFunc

func WatchObj(path string, callback func(e *clientv3.Event)) {
	ctx, cancel := context.WithCancel(context.TODO())
	ch := client.Watch(ctx, path)
	watchMap[path] = cancel
	go watching(ctx, ch, callback)
}

func CancelWatch(path string) {
	cancel := watchMap[path]
	if cancel == nil {
		return
	}
	cancel()
	delete(watchMap, path)
}

func watching(ctx context.Context, watchChan clientv3.WatchChan, callback func(e *clientv3.Event)) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-watchChan:
			for _, event := range resp.Events {
				callback(event)
			}
		}
	}
}
