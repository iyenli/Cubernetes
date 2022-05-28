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

func GetIngress(ctx *gin.Context) {
	getObj(ctx, object.IngressEtcdPrefix+ctx.Param("uid"))
}

func GetIngresses(ctx *gin.Context) {
	getObjs(ctx, object.IngressEtcdPrefix)
}

func PostIngress(ctx *gin.Context) {
	ingress := object.Ingress{}
	err := ctx.BindJSON(&ingress)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if ingress.Name == "" {
		utils.BadRequest(ctx)
		return
	}
	ingress.UID = uuid.New().String()
	buf, _ := json.Marshal(ingress)
	err = etcdrw.PutObj(object.IngressEtcdPrefix+ingress.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, ingress)
}

func PutIngress(ctx *gin.Context) {
	newIng := object.Ingress{}
	err := ctx.BindJSON(&newIng)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newIng.UID != ctx.Param("uid") {
		utils.BadRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.IngressEtcdPrefix + newIng.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newIng)
	err = etcdrw.PutObj(object.IngressEtcdPrefix+newIng.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelIngress(ctx *gin.Context) {
	delObj(ctx, object.IngressEtcdPrefix+ctx.Param("uid"))
}

func SelectIngresses(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.IngressEtcdPrefix)
		return
	}

	selectObjs(ctx, object.IngressEtcdPrefix, func(str []byte) bool {
		var ingress object.Ingress
		err = json.Unmarshal(str, &ingress)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := ingress.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}
