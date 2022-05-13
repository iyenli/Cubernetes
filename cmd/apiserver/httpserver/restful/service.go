package restful

import (
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
	getObj(ctx, "/apis/service/"+ctx.Param("uid"))
}

func GetServices(ctx *gin.Context) {
	getObjs(ctx, "/apis/service/")
}

func PostService(ctx *gin.Context) {
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
	service, err = ClusterIPAllocator.AllocateClusterIP(&service)
	if err != nil {
		serverError(ctx)
		return
	}

	buf, _ := json.Marshal(service)
	err = etcdrw.PutObj("/apis/service/"+service.UID, string(buf))
	if err != nil {
		serverError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, service)
}

func PutService(ctx *gin.Context) {
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

func DelService(ctx *gin.Context) {
	delObj(ctx, "/apis/service/"+ctx.Param("uid"))
}

func SelectServices(ctx *gin.Context) {
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
