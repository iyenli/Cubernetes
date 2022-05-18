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

func GetNode(ctx *gin.Context) {
	getObj(ctx, object.NodeEtcdPrefix+ctx.Param("uid"))
}

func GetNodes(ctx *gin.Context) {
	getObjs(ctx, object.NodeEtcdPrefix)
}

func PostNode(ctx *gin.Context) {
	var node object.Node
	err := ctx.BindJSON(&node)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}
	if node.Name == "" {
		utils.BadRequest(ctx)
		return
	}

	node.UID = uuid.New().String()
	buf, _ := json.Marshal(node)
	err = etcdrw.PutObj(object.NodeEtcdPrefix+node.UID, string(buf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, node)
}

func PutNode(ctx *gin.Context) {
	newNode := object.Node{}
	err := ctx.BindJSON(&newNode)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if newNode.UID != ctx.Param("uid") {
		utils.BadRequest(ctx)
		return
	}

	oldBuf, err := etcdrw.GetObj(object.NodeEtcdPrefix + newNode.UID)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	newBuf, _ := json.Marshal(newNode)
	err = etcdrw.PutObj(object.NodeEtcdPrefix+newNode.UID, string(newBuf))
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(newBuf))
}

func DelNode(ctx *gin.Context) {
	delObj(ctx, object.NodeEtcdPrefix+ctx.Param("uid"))
}

func SelectNodes(ctx *gin.Context) {
	var selectors map[string]string
	err := ctx.BindJSON(&selectors)
	if err != nil {
		utils.ParseFail(ctx)
		return
	}

	if len(selectors) == 0 {
		getObjs(ctx, object.NodeEtcdPrefix)
		return
	}

	selectObjs(ctx, object.NodeEtcdPrefix, func(str []byte) bool {
		var node object.Node
		err = json.Unmarshal(str, &node)
		if err != nil {
			return false
		}
		for key, val := range selectors {
			v := node.Labels[key]
			if v != val {
				return false
			}
		}
		return true
	})
}
