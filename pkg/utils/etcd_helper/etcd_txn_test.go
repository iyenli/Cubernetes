package etcd_helper

// FIX
//func TestCreateTxn(t *testing.T) {
//	ctx := ETCDContext{client: newETCDClient()}
//	defer closeETCDClient(ctx.client)
//
//	txnResponse, err := createTxn().If(
//		clientv3.Compare(clientv3.Value("test1"), "<", "test2")).
//		Then(clientv3.OpPut("test1", "test1")).
//		Else(clientv3.OpPut("test2", "test2")).
//		Commit()
//
//	if err != nil {
//		log.Panicln("Error: txn commit error")
//		return
//	}
//	assert.Equal(t, true, txnResponse.Succeeded)
//	allPods, err := ctx.client.KV.Get(context.TODO(), "test1")
//	if err != nil {
//		log.Panicln("Error: get all pod error")
//		return
//	}
//
//	assert.Equal(t, 1, len(allPods.Kvs))
//	assert.Equal(t, []byte("test1"), allPods.Kvs[0].Value)
//}
