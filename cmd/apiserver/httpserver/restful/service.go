package restful

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	"Cubernetes/pkg/cubenetwork/servicenetwork"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var ClusterIPAllocator *servicenetwork.ClusterIPAllocator

func GetService(ctx *gin.Context) {
	getObj(ctx, object.ServiceEtcdPrefix+ctx.Param("uid"))
}

func GetServices(ctx *gin.Context) {
	getObjs(ctx, object.ServiceEtcdPrefix)
}

func PostService(ctx *gin.Context) {
	service := object.Service{}
	err := ctx.BindJSON(&service)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if service.Name == "" {
		utils.BadRequest(ctx)
		return
	}
	service.UID = uuid.New().String()
	service, err = ClusterIPAllocator.AllocateClusterIP(&service)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	buf, _ := json.Marshal(service)
	err = etcdrw.PutObj(object.ServiceEtcdPrefix+service.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, service)
}

func PutService(ctx *gin.Context) {
	newService := object.Service{}
	err := ctx.BindJSON(&newService)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newService.UID != ctx.Param("uid") {
		utils.BadRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.ServiceEtcdPrefix + newService.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newService)
	err = etcdrw.PutObj(object.ServiceEtcdPrefix+newService.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelService(ctx *gin.Context) {
	delObj(ctx, object.ServiceEtcdPrefix+ctx.Param("uid"))
}

func SelectServices(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.ServiceEtcdPrefix)
		return
	}

	selectObjs(ctx, object.ServiceEtcdPrefix, func(str []byte) bool {
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
