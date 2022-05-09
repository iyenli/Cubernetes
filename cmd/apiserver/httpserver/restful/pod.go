package restful

import (
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetPod(ctx *gin.Context) {
	getObj(ctx, "/apis/pod/"+ctx.Param("uid"))
}

func GetPods(ctx *gin.Context) {
	getObjs(ctx, "/apis/pod/")
}

func PostPod(ctx *gin.Context) {
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

func PutPod(ctx *gin.Context) {
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

func DelPod(ctx *gin.Context) {
	delObj(ctx, "/apis/pod/"+ctx.Param("uid"))
}

func SelectPods(ctx *gin.Context) {
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

func UpdatePodStatus(ctx *gin.Context) {
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
