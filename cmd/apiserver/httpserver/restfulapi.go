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
	{http.MethodPost, "/apis/select/pods", selectPods},
	{http.MethodPut, "/apis/pod/status/:uid", updatePodStatus},

	{http.MethodGet, "/apis/service/:uid", getService},
	{http.MethodGet, "/apis/services", getServices},
	{http.MethodPost, "/apis/service", postService},
	{http.MethodPut, "/apis/service/:uid", putService},
	{http.MethodDelete, "/apis/service/:uid", delService},
	{http.MethodPost, "/apis/select/services", selectServices},

	{http.MethodGet, "/apis/replicaSet/:uid", getReplicaSet},
	{http.MethodGet, "/apis/replicaSets", getReplicaSets},
	{http.MethodPost, "/apis/replicaSet", postReplicaSet},
	{http.MethodPut, "/apis/replicaSet/:uid", putReplicaSet},
	{http.MethodDelete, "/apis/replicaSet/:uid", delReplicaSet},
	{http.MethodPost, "/apis/select/replicaSets", selectReplicaSets},

	{http.MethodPost, "/apis/node", postNode},
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
		ctx.Header("Content-Type", "application/json")
		ctx.String(http.StatusOK, "[]")
		return
	}

	objs := "["
	for _, str := range buf {
		objs += string(str) + ","
	}
	if len(objs) > 1 {
		objs = objs[:len(objs)-1]
	}
	objs += "]"

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, objs)
}

func selectObjs(ctx *gin.Context, prefix string, match func([]byte) bool) {
	buf, err := etcdrw.GetObjs(prefix)
	if err != nil {
		serverError(ctx)
		return
	}
	if buf == nil {
		ctx.Header("Content-Type", "application/json")
		ctx.String(http.StatusOK, "[]")
		return
	}

	objs := "["
	for _, str := range buf {
		if match(str) {
			objs += string(str) + ","
		}
	}
	if len(objs) > 1 {
		objs = objs[:len(objs)-1]
	}
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
	pod.UID = uuid.New().String()
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

	oldBuf, err := etcdrw.GetObj("/apis/pod/" + newPod.UID)
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

func selectPods(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		parseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, "/apis/pod/")
		return
	}

	selectObjs(ctx, "/apis/pod/", func(str []byte) bool {
		var pod object.Pod
		err = json.Unmarshal(str, &pod)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := pod.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}

func updatePodStatus(ctx *gin.Context) {
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

	buf, err := etcdrw.GetObj("/apis/pod/" + newPod.UID)
	if err != nil {
		serverError(ctx)
		return
	}
	if buf == nil {
		notFound(ctx)
		return
	}

	var pod object.Pod
	err = json.Unmarshal(buf, &pod)
	if err != nil {
		serverError(ctx)
		return
	}

	pod.Status = newPod.Status
	newBuf, _ := json.Marshal(pod)
	err = etcdrw.PutObj("/apis/pod/"+newPod.UID, string(newBuf))
	if err != nil {
		serverError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
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
	service.UID = uuid.New().String()
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

	oldBuf, err := etcdrw.GetObj("/apis/service/" + newService.UID)
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

func selectServices(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		parseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, "/apis/service/")
		return
	}

	selectObjs(ctx, "/apis/service/", func(str []byte) bool {
		var service object.Service
		err = json.Unmarshal(str, &service)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := service.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}

func getReplicaSet(ctx *gin.Context) {
	getObj(ctx, "/apis/replicaSet/"+ctx.Param("uid"))
}

func getReplicaSets(ctx *gin.Context) {
	getObjs(ctx, "/apis/replicaSet/")
}

func postReplicaSet(ctx *gin.Context) {
	rs := object.ReplicaSet{}
	err := ctx.BindJSON(&rs)
	if err != nil {
		parseFail(ctx)
		return
	}
	if rs.Name == "" {
		badRequest(ctx)
		return
	}
	rs.UID = uuid.New().String()
	buf, _ := json.Marshal(rs)
	err = etcdrw.PutObj("/apis/replicaSet/"+rs.UID, string(buf))
	if err != nil {
		serverError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, rs)
}

func putReplicaSet(ctx *gin.Context) {
	newRs := object.ReplicaSet{}
	err := ctx.BindJSON(&newRs)
	if err != nil {
		parseFail(ctx)
		return
	}

	if newRs.UID != ctx.Param("uid") {
		badRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj("/apis/replicaSet/" + newRs.UID)
	if err != nil {
		serverError(ctx)
		return
	}
	if oldBuf == nil {
		notFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newRs)
	err = etcdrw.PutObj("/apis/replicaSet/"+newRs.UID, string(newBuf))
	if err != nil {
		serverError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func delReplicaSet(ctx *gin.Context) {
	delObj(ctx, "/apis/replicaSet/"+ctx.Param("uid"))
}

func selectReplicaSets(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		parseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, "/apis/replicaSet/")
		return
	}

	selectObjs(ctx, "/apis/replicaSet/", func(str []byte) bool {
		var rs object.ReplicaSet
		err = json.Unmarshal(str, &rs)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := rs.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}

func postNode(ctx *gin.Context) {
	var node object.Node
	err := ctx.BindJSON(&node)
	if err != nil {
		parseFail(ctx)
		return
	}
	if node.Name == "" {
		badRequest(ctx)
		return
	}

	node.UID = uuid.New().String()
	buf, _ := json.Marshal(node)
	err = etcdrw.PutObj("/apis/node/"+node.UID, string(buf))
	if err != nil {
		serverError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, node)
}
