package httpserver

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
)

var watchList = []Handler{
	{http.MethodPost, "/apis/watch/pod/:uid", watchPod},
	{http.MethodPost, "/apis/watch/pods", watchPods},

	{http.MethodPost, "/apis/watch/service/:uid", watchService},
	{http.MethodPost, "/apis/watch/services", watchServices},

	{http.MethodPost, "/apis/watch/replicaSet/:uid", watchReplicaSet},
	{http.MethodPost, "/apis/watch/replicaSets", watchReplicaSets},

	{http.MethodPost, "/apis/watch/node/:uid", watchNode},
	{http.MethodPost, "/apis/watch/nodes", watchNodes},

	{http.MethodPost, "/apis/watch/dns/:uid", watchDns},
	{http.MethodPost, "/apis/watch/dnses", watchDnses},

	{http.MethodPost, "/apis/watch/autoScaler/:uid", watchAutoScaler},
	{http.MethodPost, "/apis/watch/autoScalers", watchAutoScalers},

	{http.MethodPost, "/apis/watch/gpuJob/:uid", watchGpuJob},
	{http.MethodPost, "/apis/watch/gpuJobs", watchGpuJobs},

	{http.MethodPost, "/apis/watch/action/:uid", watchAction},
	{http.MethodPost, "/apis/watch/actions", watchActions},

	{http.MethodPost, "/apis/watch/actor/:uid", watchActor},
	{http.MethodPost, "/apis/watch/actors", watchActors},

	{http.MethodPost, "/apis/watch/ingress/:uid", watchIngress},
	{http.MethodPost, "/apis/watch/ingresses", watchIngresses},
}

func handleEvent(ctx *gin.Context, e *clientv3.Event) {
	log.Println("watched event, telling client...")
	var objEvent watchobj.ObjEvent
	switch e.Type {
	case mvccpb.PUT:
		objEvent.EType = watchobj.EVENT_PUT
	case mvccpb.DELETE:
		objEvent.EType = watchobj.EVENT_DELETE
	}
	objEvent.Path = string(e.Kv.Key)
	objEvent.Object = string(e.Kv.Value)
	buf, _ := json.Marshal(objEvent)
	buf = append(buf, watchobj.MSG_DELIM)
	_, err := ctx.Writer.Write(buf)
	if err != nil {
		log.Println("fail to write to http client, error: ", err)
		return
	}
	ctx.Writer.Flush()
}

func postWatch(ctx *gin.Context, path string, withPrefix bool) {
	c, cancel := context.WithCancel(context.TODO())

	var watchChan clientv3.WatchChan
	if withPrefix {
		watchChan = etcdrw.WatchObjs(c, path)
	} else {
		watchChan = etcdrw.WatchObj(c, path)
	}

	buf := []byte(watchobj.WATCH_CONFIRM)
	buf = append(buf, watchobj.MSG_DELIM)
	_, err := ctx.Writer.Write(buf)
	if err != nil {
		log.Println("fail to write to http client, error: ", err)
		cancel()
		return
	}
	ctx.Writer.Flush()

	for {
		select {
		case <-ctx.Request.Context().Done():
			log.Println("connection closed, canceling watch...")
			cancel()
			return
		case resp := <-watchChan:
			for _, event := range resp.Events {
				handleEvent(ctx, event)
			}
		}
	}
}

func watchPod(ctx *gin.Context) {
	postWatch(ctx, object.PodEtcdPrefix+ctx.Param("uid"), false)
}

func watchPods(ctx *gin.Context) {
	postWatch(ctx, object.PodEtcdPrefix, true)
}

func watchService(ctx *gin.Context) {
	postWatch(ctx, object.ServiceEtcdPrefix+ctx.Param("uid"), false)
}

func watchServices(ctx *gin.Context) {
	postWatch(ctx, object.ServiceEtcdPrefix, true)
}

func watchReplicaSet(ctx *gin.Context) {
	postWatch(ctx, object.ReplicaSetEtcdPrefix+ctx.Param("uid"), false)
}

func watchReplicaSets(ctx *gin.Context) {
	postWatch(ctx, object.ReplicaSetEtcdPrefix, true)
}

func watchNode(ctx *gin.Context) {
	postWatch(ctx, object.NodeEtcdPrefix+ctx.Param("uid"), false)
}

func watchNodes(ctx *gin.Context) {
	postWatch(ctx, object.NodeEtcdPrefix, true)
}

func watchDns(ctx *gin.Context) {
	postWatch(ctx, object.DnsEtcdPrefix+ctx.Param("uid"), false)
}

func watchDnses(ctx *gin.Context) {
	postWatch(ctx, object.DnsEtcdPrefix, true)
}

func watchAutoScaler(ctx *gin.Context) {
	postWatch(ctx, object.AutoScalerEtcdPrefix+ctx.Param("uid"), false)
}

func watchAutoScalers(ctx *gin.Context) {
	postWatch(ctx, object.AutoScalerEtcdPrefix, true)
}

func watchGpuJob(ctx *gin.Context) {
	postWatch(ctx, object.GpuJobEtcdPrefix+ctx.Param("uid"), false)
}

func watchGpuJobs(ctx *gin.Context) {
	postWatch(ctx, object.GpuJobEtcdPrefix, true)
}

func watchAction(ctx *gin.Context) {
	postWatch(ctx, object.ActionEtcdPrefix+ctx.Param("uid"), false)
}

func watchActions(ctx *gin.Context) {
	postWatch(ctx, object.ActionEtcdPrefix, true)
}

func watchActor(ctx *gin.Context) {
	postWatch(ctx, object.ActorEtcdPrefix+ctx.Param("uid"), false)
}

func watchActors(ctx *gin.Context) {
	postWatch(ctx, object.ActorEtcdPrefix, true)
}

func watchIngress(ctx *gin.Context) {
	postWatch(ctx, object.IngressEtcdPrefix+ctx.Param("uid"), false)
}

func watchIngresses(ctx *gin.Context) {
	postWatch(ctx, object.IngressEtcdPrefix, true)
}
