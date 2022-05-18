package restful

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetPod(ctx *gin.Context) {
	getObj(ctx, object.PodEtcdPrefix+ctx.Param("uid"))
}

func GetPods(ctx *gin.Context) {
	getObjs(ctx, object.PodEtcdPrefix)
}

func PostPod(ctx *gin.Context) {
	pod := object.Pod{}
	err := ctx.BindJSON(&pod)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if pod.Name == "" {
		utils.BadRequest(ctx)
		return
	}
	pod.UID = uuid.New().String()
	buf, _ := json.Marshal(pod)
	err = etcdrw.PutObj(object.PodEtcdPrefix+pod.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, pod)
}

func PutPod(ctx *gin.Context) {
	newPod := object.Pod{}
	err := ctx.BindJSON(&newPod)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newPod.UID != ctx.Param("uid") {
		utils.BadRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.PodEtcdPrefix + newPod.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newPod)
	err = etcdrw.PutObj(object.PodEtcdPrefix+newPod.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelPod(ctx *gin.Context) {
	delObj(ctx, object.PodEtcdPrefix+ctx.Param("uid"))
}

func SelectPods(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.PodEtcdPrefix)
		return
	}

	selectObjs(ctx, object.PodEtcdPrefix, func(str []byte) bool {
		var pod object.Pod
		err = json.Unmarshal(str, &pod)
		if err != nil {
			return false
		}

		return object.MatchLabelSelector(selectors, pod.Labels)
	})
}

func UpdatePodStatus(ctx *gin.Context) {
	newPod := object.Pod{}
	err := ctx.BindJSON(&newPod)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newPod.UID != ctx.Param("uid") {
		utils.BadRequest(ctx)
		return
	}

	buf, err := etcdrw.GetObj(object.PodEtcdPrefix + newPod.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if buf == nil {
		utils.NotFound(ctx)
		return
	}

	var pod object.Pod
	err = json.Unmarshal(buf, &pod)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	pod.Status = newPod.Status
	newBuf, _ := json.Marshal(pod)
	err = etcdrw.PutObj(object.PodEtcdPrefix+newPod.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}
