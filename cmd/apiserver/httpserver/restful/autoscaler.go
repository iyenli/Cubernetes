package restful

import (
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetAutoScaler(ctx *gin.Context) {
	getObj(ctx, object.AutoScalerEtcdPrefix+ctx.Param("uid"))
}

func GetAutoScalers(ctx *gin.Context) {
	getObjs(ctx, object.AutoScalerEtcdPrefix)
}

func PostAutoScaler(ctx *gin.Context) {
	as := object.AutoScaler{}
	err := ctx.BindJSON(&as)
	if err != nil {
		parseFail(ctx)
		return
	}
	if as.Name == "" {
		badRequest(ctx)
		return
	}
	as.UID = uuid.New().String()
	buf, _ := json.Marshal(as)
	err = etcdrw.PutObj(object.AutoScalerEtcdPrefix+as.UID, string(buf))
	if err != nil {
		serverError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, as)
}

func PutAutoScaler(ctx *gin.Context) {
	newAs := object.AutoScaler{}
	err := ctx.BindJSON(&newAs)
	if err != nil {
		parseFail(ctx)
		return
	}

	if newAs.UID != ctx.Param("uid") {
		badRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.AutoScalerEtcdPrefix + newAs.UID)
	if err != nil {
		serverError(ctx)
		return
	}
	if oldBuf == nil {
		notFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newAs)
	err = etcdrw.PutObj(object.AutoScalerEtcdPrefix+newAs.UID, string(newBuf))
	if err != nil {
		serverError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelAutoScaler(ctx *gin.Context) {
	delObj(ctx, object.AutoScalerEtcdPrefix+ctx.Param("uid"))
}

func SelectAutoScalers(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		parseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.AutoScalerEtcdPrefix)
		return
	}

	selectObjs(ctx, object.AutoScalerEtcdPrefix, func(str []byte) bool {
		var as object.AutoScaler
		err = json.Unmarshal(str, &as)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := as.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}
