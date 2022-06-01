package restful

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/dag"
	"Cubernetes/pkg/utils/etcdrw"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
)

func GetAction(ctx *gin.Context) {
	getObj(ctx, object.ActionEtcdPrefix+ctx.Param("uid"))
}

func GetActions(ctx *gin.Context) {
	getObjs(ctx, object.ActionEtcdPrefix)
}

func PostAction(ctx *gin.Context) {
	newAction := object.Action{}
	err := ctx.BindJSON(&newAction)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if newAction.Name == "" {
		utils.BadRequest(ctx)
		return
	}

	actions, err := etcdrw.GetObjs(object.ActionEtcdPrefix)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	existed := false
	nodes := make(map[string][]string)
	for _, buf := range actions {
		var action object.Action
		err = json.Unmarshal(buf, &action)
		if err != nil {
			utils.ServerError(ctx)
			return
		}
		if newAction.Name == action.Name {
			// Action existed, update the action
			newAction.UID = action.UID
			newAction.Status = action.Status
			existed = true
			continue
		}
		nodes[action.Name] = action.Spec.InvokeActions
	}
	nodes[newAction.Name] = newAction.Spec.InvokeActions
	containCircle, cycle := dag.CheckLoop(nodes)
	if containCircle {
		ret := fmt.Sprintf("new action will form a cycle: %v", cycle)
		ctx.String(http.StatusBadRequest, ret)
		return
	}

	if !existed {
		newAction.UID = uuid.New().String()
	}

	buf, _ := json.Marshal(newAction)
	err = etcdrw.PutObj(object.ActionEtcdPrefix+newAction.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, newAction)
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

	actions, err := etcdrw.GetObjs(object.ActionEtcdPrefix)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	existed := false
	nodes := make(map[string][]string)
	for _, buf := range actions {
		var action object.Action
		err = json.Unmarshal(buf, &action)
		if err != nil {
			utils.ServerError(ctx)
			return
		}
		if newAct.UID == action.UID {
			existed = true
			continue
		}
		nodes[action.Name] = action.Spec.InvokeActions
	}

	if !existed {
		utils.NotFound(ctx)
		return
	}

	nodes[newAct.Name] = newAct.Spec.InvokeActions
	containCircle, cycle := dag.CheckLoop(nodes)
	if containCircle {
		ret := fmt.Sprintf("new action will form a cycle: %v", cycle)
		ctx.String(http.StatusBadRequest, ret)
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
	key := object.ActionEtcdPrefix + ctx.Param("uid")
	oldBuf, err := etcdrw.GetObj(key)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	err = etcdrw.DelObj(key)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.String(http.StatusOK, "deleted")

	var action object.Action
	err = json.Unmarshal(oldBuf, &action)
	if err != nil || action.Name == "" {
		return
	}

	filename := path.Join(cubeconfig.ActionFileDir, action.Spec.ScriptUID)
	_ = os.RemoveAll(filename + ".py")
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

type workflow struct {
	Ingresses []ingress `json:"ingresses"`
	Actions   []action  `json:"actions"`
}

type ingress struct {
	Src  string `json:"s"`
	Dest string `json:"d"`
	Path string `json:"p"`
}

type action struct {
	Src  string   `json:"s"`
	Dest []string `json:"d"`
}

func GetWorkflow(ctx *gin.Context) {
	actions := make([]action, 0)
	buf, err := etcdrw.GetObjs(object.ActionEtcdPrefix)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	for _, actStr := range buf {
		var act object.Action
		err = json.Unmarshal(actStr, &act)
		if err != nil {
			utils.ServerError(ctx)
			return
		}
		actions = append(actions, action{Src: act.Name, Dest: act.Spec.InvokeActions})
	}

	ingresses := make([]ingress, 0)
	buf, err = etcdrw.GetObjs(object.IngressEtcdPrefix)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	for _, igsStr := range buf {
		var igs object.Ingress
		err = json.Unmarshal(igsStr, &igs)
		if err != nil {
			utils.ServerError(ctx)
			return
		}
		ingresses = append(ingresses, ingress{Src: igs.Name, Dest: igs.Spec.InvokeAction, Path: igs.Spec.TriggerPath})
	}

	wf := workflow{
		Ingresses: ingresses,
		Actions:   actions,
	}

	ctx.JSON(http.StatusOK, wf)
}
