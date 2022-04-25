package httpserver

import (
	"Cubernetes/pkg/utils/etcdrw"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/clientv3"
	"log"
	"net/http"
)

var watchList = []Handler{
	{http.MethodPost, "/apis/watch/pod/:uid", watchPod},
}

func handleEvent(ctx *gin.Context, e *clientv3.Event) {
	flusher, _ := ctx.Writer.(http.Flusher)
	log.Println("watched event, telling client...")
	_, err := fmt.Fprint(ctx.Writer, string(e.Kv.Value))
	if err != nil {
		log.Println("fail to write to http client, error: ", err)
		return
	}
	flusher.Flush()
}

func postWatch(ctx *gin.Context, path string) {
	c, cancel := context.WithCancel(context.TODO())
	watchChan := etcdrw.WatchObj(c, path)
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
	postWatch(ctx, "/apis/pod/"+ctx.Param("uid"))
}
