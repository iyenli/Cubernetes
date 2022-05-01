package etcdrw

import (
	"context"
	"go.etcd.io/etcd/client/v3"
)

func WatchObj(ctx context.Context, path string) clientv3.WatchChan {
	ch := client.Watch(ctx, path)
	return ch
}

func WatchObjs(ctx context.Context, prefix string) clientv3.WatchChan {
	ch := client.Watch(ctx, prefix, clientv3.WithPrefix())
	return ch
}
