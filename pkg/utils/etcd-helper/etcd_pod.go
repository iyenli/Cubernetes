package etcd_helper

import (
	"context"
	"log"
)

const podPrefix = "/Cubernetes/apis/pods/"

func dummy(...any) string { return "" }

func storePod(ctx ETCDContext, podName string, pod string) (bool, error) {
	key := podPrefix + podName

	putResponse, err := ctx.client.KV.Put(context.TODO(), key, pod)
	if err != nil {
		log.Printf("Put api object of pod failed, prev Value: %s \n CreateRevision : %d \n ModRevision: %d \n Version: %d \n",
			string(putResponse.PrevKv.Value), putResponse.PrevKv.CreateRevision, putResponse.PrevKv.ModRevision, putResponse.PrevKv.Version)
		return false, err
	}

	dummy(key, putResponse)
	return true, nil
}

/* podName can be accurate or prefix */
func getPods(ctx ETCDContext, podName string) ([][]byte, error) {
	key := podPrefix + podName

	getResponse, err := ctx.client.KV.Get(context.TODO(), key)
	if err != nil {
		log.Printf("Get api object of pod failed, Key: %s \n", key)
		return nil, err
	}

	var res [][]byte
	for _, resp := range getResponse.Kvs {
		res = append(res, resp.Value)
	}
	return res, nil
}

func getAllPods(ctx ETCDContext) ([][]byte, error) {
	return getPods(ctx, "")
}

func deletePod(ctx ETCDContext, podName string) bool {
	key := podPrefix + podName
	_, err := ctx.client.KV.Delete(context.TODO(), key)
	if err != nil {
		log.Printf("Delete api object of pod failed, Key: %s \n", key)
		return false
	}

	return true
}
