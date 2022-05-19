package restful

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	"Cubernetes/pkg/utils/etcdrw"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getObj(ctx *gin.Context, path string) {
	buf, err := etcdrw.GetObj(path)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if buf == nil {
		utils.NotFound(ctx)
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, string(buf))
}

func getObjs(ctx *gin.Context, prefix string) {
	buf, err := etcdrw.GetObjs(prefix)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if buf == nil {
		ctx.Header("Content-Type", "application/json")
		ctx.String(http.StatusOK, "[]")
		return
	}

	objs := "["
	for _, str := range buf {
		objs += string(str) + ","
	}
	if len(objs) > 1 {
		objs = objs[:len(objs)-1]
	}
	objs += "]"

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, objs)
}

func selectObjs(ctx *gin.Context, prefix string, match func([]byte) bool) {
	buf, err := etcdrw.GetObjs(prefix)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if buf == nil {
		ctx.Header("Content-Type", "application/json")
		ctx.String(http.StatusOK, "[]")
		return
	}

	objs := "["
	for _, str := range buf {
		if match(str) {
			objs += string(str) + ","
		}
	}
	if len(objs) > 1 {
		objs = objs[:len(objs)-1]
	}
	objs += "]"

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, objs)
}

func delObj(ctx *gin.Context, path string) {
	oldBuf, err := etcdrw.GetObj(path)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	if oldBuf == nil {
		utils.NotFound(ctx)
		return
	}

	err = etcdrw.DelObj(path)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	ctx.String(http.StatusOK, "deleted")
}
