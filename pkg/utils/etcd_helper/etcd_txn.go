package etcd_helper

import "go.etcd.io/etcd/clientv3"

func createTxn(ctx *ETCDContext) clientv3.Txn {
	return ctx.client.KV.Txn(ctx.client.Ctx())
}
