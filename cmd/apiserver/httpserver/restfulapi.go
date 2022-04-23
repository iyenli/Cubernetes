package httpserver

import (
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

var restfulList = []Handler{
	{http.MethodGet, "/apis/pod/:name", getPod},
	{http.MethodPost, "/apis/pod", postPod},
}

func getPod(ctx *gin.Context) {
	podStr, err := etcdrw.GetObj("/apis/pod/" + ctx.Param("name"))
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, podStr)
}

func postPod(ctx *gin.Context) {
	pod := object.Pod{}
	err := ctx.BindJSON(&pod)
	if err != nil {
		return
	}
	name := pod.Name
	if name != "" {
		podstr, _ := json.Marshal(pod)
		etcdrw.PutObj("/apis/pod/"+pod.Name, string(podstr))
	}
}
