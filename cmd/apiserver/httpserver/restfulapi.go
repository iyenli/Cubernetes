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

	{http.MethodGet, "/apis/service/:uid", getService},
	{http.MethodGet, "/apis/services", getServices},
	{http.MethodPost, "/apis/service", postService},
	{http.MethodPut, "/apis/service/:uid", putService},
	{http.MethodDelete, "/apis/service/:uid", delService},
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

func getObj(ctx *gin.Context, path string) {
	buf, err := etcdrw.GetObj(path)
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

func getObjs(ctx *gin.Context, prefix string) {
	buf, err := etcdrw.GetObjs(prefix)
	if err != nil {
		serverError(ctx)
		return
	}
	if buf == nil {
		notFound(ctx)
		return
	}

	objs := "["
	for _, str := range buf {
		objs += string(str) + ","
	}
	objs = objs[:len(objs)-1]
	objs += "]"

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, objs)
}

func delObj(ctx *gin.Context, path string) {
	oldBuf, err := etcdrw.GetObj(path)
	if err != nil {
		serverError(ctx)
		return
	}
	if oldBuf == nil {
		notFound(ctx)
		return
	}

	err = etcdrw.DelObj(path)
	if err != nil {
		serverError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, "deleted")
}

func getPod(ctx *gin.Context) {
	getObj(ctx, "/apis/pod/"+ctx.Param("uid"))
}

func getPods(ctx *gin.Context) {
	getObjs(ctx, "/apis/pod/")
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
	delObj(ctx, "/apis/pod/"+ctx.Param("uid"))
}

func getService(ctx *gin.Context) {
	getObj(ctx, "/apis/service/"+ctx.Param("uid"))
}

func getServices(ctx *gin.Context) {
	getObjs(ctx, "/apis/service/")
}

func postService(ctx *gin.Context) {
	service := object.Service{}
	err := ctx.BindJSON(&service)
	if err != nil {
		parseFail(ctx)
		return
	}
	if service.Name == "" {
		badRequest(ctx)
		return
	}
	service.UID = service.Name + ":" + uuid.New().String()
	buf, _ := json.Marshal(service)
	err = etcdrw.PutObj("/apis/service/"+service.UID, string(buf))
	if err != nil {
		serverError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, service)
}

func putService(ctx *gin.Context) {
	newService := object.Service{}
	err := ctx.BindJSON(&newService)
	if err != nil {
		parseFail(ctx)
		return
	}

	if newService.UID != ctx.Param("uid") {
		badRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj("/apis/service/" + ctx.Param("uid"))
	if err != nil {
		serverError(ctx)
		return
	}
	if oldBuf == nil {
		notFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newService)
	err = etcdrw.PutObj("/apis/service/"+newService.UID, string(newBuf))
	if err != nil {
		serverError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func delService(ctx *gin.Context) {
	delObj(ctx, "/apis/service/"+ctx.Param("uid"))
}
