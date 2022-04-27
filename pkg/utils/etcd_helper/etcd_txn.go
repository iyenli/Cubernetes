package etcd_helper

import "go.etcd.io/etcd/clientv3"

func CreateTxn(ctx *ETCDContext) clientv3.Txn {
	return ctx.Client.KV.Txn(ctx.Client.Ctx())
}
