package httpserver

import (
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcd_helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Method     string
	Path       string
	HandleFunc func(ctx *gin.Context)
}

var handlerList = [...]Handler{
	{http.MethodGet, "/api/pod/:name", getPod},
	{http.MethodPost, "/api/pod", postPod},
}

func getPod(ctx *gin.Context) {
	podStr, err := etcd_helper.GetPods(&etcd, ctx.Param("name"))
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
}
