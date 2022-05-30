package file

import (
	cubeconfig "Cubernetes/config"
	"github.com/gin-gonic/gin"
	"path"
)

func GetActionFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.ActionFileDir, ctx.Param("uid")+".py")
	GetFile(ctx, filename)
}

func PostActionFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.ActionFileDir, ctx.Param("uid")+".py")
	PostFile(ctx, filename)
}
