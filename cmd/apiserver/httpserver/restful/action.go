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

func GetAction(ctx *gin.Context) {
	getObj(ctx, object.ActionEtcdPrefix+ctx.Param("uid"))
}

func GetActions(ctx *gin.Context) {
	getObjs(ctx, object.ActionEtcdPrefix)
}

func PostAction(ctx *gin.Context) {
	action := object.Action{}
	err := ctx.BindJSON(&action)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if action.Name == "" {
		utils.BadRequest(ctx)
		return
	}

	acts, err := etcdrw.GetObjs(object.ActionEtcdPrefix)
	for _, buf := range acts {
		var act object.Action
		err = json.Unmarshal(buf, &act)
		if err != nil {
			utils.ServerError(ctx)
			return
		}
		if act.Name == action.Name {
			ctx.String(http.StatusBadRequest, "bad request: Action %s existed", action.Name)
			return
		}
	}

	action.UID = uuid.New().String()
	action.Status = &object.ActionStatus{
		Phase: object.ActionCreating,
	}

	buf, _ := json.Marshal(action)
	err = etcdrw.PutObj(object.ActionEtcdPrefix+action.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, action)
}

func PutAction(ctx *gin.Context) {
	newAct := object.Action{}
	err := ctx.BindJSON(&newAct)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newAct.UID != ctx.Param("uid") {
		utils.BadRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.ActionEtcdPrefix + newAct.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newAct)
	err = etcdrw.PutObj(object.ActionEtcdPrefix+newAct.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelAction(ctx *gin.Context) {
	delObj(ctx, object.ActionEtcdPrefix+ctx.Param("uid"))
}

func SelectActions(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.ActionEtcdPrefix)
		return
	}

	selectObjs(ctx, object.ActionEtcdPrefix, func(str []byte) bool {
		var action object.Action
		err = json.Unmarshal(str, &action)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := action.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}
