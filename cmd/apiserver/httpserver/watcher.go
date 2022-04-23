package httpserver

import (
	"Cubernetes/pkg/utils/etcdrw"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/clientv3"
	"net/http"
)

var watchList = []Handler{
	{http.MethodPost, "/apis/watch/:path", postWatch},
}

func postWatch(ctx *gin.Context) {
	flusher, _ := ctx.Writer.(http.Flusher)
	etcdrw.WatchObj(ctx.Param("path"), func(event *clientv3.Event) {
		switch event.Type {
		case mvccpb.PUT:
			fmt.Println(event.Kv.Key, "changed")
			_, err := fmt.Fprintf(ctx.Writer, "watch result key: %v, value: %v\n", event.Kv.Key, event.Kv.Value)
			if err != nil {
				return
			}
			flusher.Flush()
		case mvccpb.DELETE:
			fmt.Println(event.Kv.Key, "deleted")
		}
	})
}
