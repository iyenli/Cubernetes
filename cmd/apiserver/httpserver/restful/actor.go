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

func GetActor(ctx *gin.Context) {
	getObj(ctx, object.ActorEtcdPrefix+ctx.Param("uid"))
}

func GetActors(ctx *gin.Context) {
	getObjs(ctx, object.ActorEtcdPrefix)
}

func PostActor(ctx *gin.Context) {
	actor := object.Actor{}
	err := ctx.BindJSON(&actor)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if actor.Name == "" {
		utils.BadRequest(ctx)
		return
	}
	actor.UID = uuid.New().String()
	buf, _ := json.Marshal(actor)
	err = etcdrw.PutObj(object.ActorEtcdPrefix+actor.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, actor)
}

func PutActor(ctx *gin.Context) {
	newAct := object.Actor{}
	err := ctx.BindJSON(&newAct)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newAct.UID != ctx.Param("uid") {
		utils.BadRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.ActorEtcdPrefix + newAct.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newAct)
	err = etcdrw.PutObj(object.ActorEtcdPrefix+newAct.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelActor(ctx *gin.Context) {
	delObj(ctx, object.ActorEtcdPrefix+ctx.Param("uid"))
}

func SelectActors(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.ActorEtcdPrefix)
		return
	}

	selectObjs(ctx, object.ActorEtcdPrefix, func(str []byte) bool {
		var actor object.Actor
		err = json.Unmarshal(str, &actor)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := actor.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}
