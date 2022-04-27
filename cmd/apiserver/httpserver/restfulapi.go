package httpserver

import (
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var restfulList = []Handler{
	{http.MethodGet, "/apis/pod/:uid", getPod},
	{http.MethodGet, "/apis/pods", getPods},
	{http.MethodPost, "/apis/pod", postPod},
	{http.MethodPut, "/apis/pod/:uid", putPod},
	{http.MethodDelete, "/apis/pod/:uid", delPod},
}

func parseFail(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, "fail to parse body")
}

func badRequest(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, "bad request")
}

func serverError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, "server error")
}

func notFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, "no objects found")
}

func getPod(ctx *gin.Context) {
	buf, err := etcdrw.GetObj("/apis/pod/" + ctx.Param("uid"))
	if err != nil {
		serverError(ctx)
		return
	}
	if buf == nil {
		notFound(ctx)
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(buf))
}

func getPods(ctx *gin.Context) {
	buf, err := etcdrw.GetObjs("/apis/pod")
	if err != nil {
		serverError(ctx)
		return
	}
	if buf == nil {
		notFound(ctx)
		return
	}

	pods := "["
	for _, str := range buf {
		pods += string(str) + ","
	}
	pods = pods[:len(pods)-1]
	pods += "]"

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, pods)
}

func postPod(ctx *gin.Context) {
	pod := object.Pod{}
	err := ctx.BindJSON(&pod)
	if err != nil {
		parseFail(ctx)
		return
	}
	if pod.Name == "" {
		badRequest(ctx)
		return
	}
	pod.UID = pod.Name + ":" + uuid.New().String()
	buf, _ := json.Marshal(pod)
	err = etcdrw.PutObj("/apis/pod/"+pod.UID, string(buf))
	if err != nil {
		serverError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, pod)
}

func putPod(ctx *gin.Context) {
	newPod := object.Pod{}
	err := ctx.BindJSON(&newPod)
	if err != nil {
		parseFail(ctx)
		return
	}

	if newPod.UID != ctx.Param("uid") {
		badRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj("/apis/pod/" + ctx.Param("uid"))
	if err != nil {
		serverError(ctx)
		return
	}
	if oldBuf == nil {
		notFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newPod)
	err = etcdrw.PutObj("/apis/pod/"+newPod.UID, string(newBuf))
	if err != nil {
		serverError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func delPod(ctx *gin.Context) {
	oldBuf, err := etcdrw.GetObj("/apis/pod/" + ctx.Param("uid"))
	if err != nil {
		serverError(ctx)
		return
	}
	if oldBuf == nil {
		notFound(ctx)
		return
	}

	err = etcdrw.DelObj("/apis/pod/" + ctx.Param("uid"))
	if err != nil {
		serverError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, "deleted")
}
