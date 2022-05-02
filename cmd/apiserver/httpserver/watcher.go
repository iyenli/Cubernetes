package httpserver

import (
	"Cubernetes/pkg/apiserver/watchobj"
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
	postWatch(ctx, "/apis/pod/"+ctx.Param("uid"), false)
}

func watchPods(ctx *gin.Context) {
	postWatch(ctx, "/apis/pod", true)
}

func watchService(ctx *gin.Context) {
	postWatch(ctx, "/apis/pod/"+ctx.Param("uid"), false)
}

func watchServices(ctx *gin.Context) {
	postWatch(ctx, "/apis/pod", true)
}
