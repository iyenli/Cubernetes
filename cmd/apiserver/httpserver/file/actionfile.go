package file

import (
	cubeconfig "Cubernetes/config"
	"github.com/gin-gonic/gin"
	"path"
)

func GetActionFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.ActionFileDir, ctx.Param("uid")+".tar.gz")
	GetFile(ctx, filename)
}

func PostActionFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.ActionFileDir, ctx.Param("uid")+".tar.gz")
	PostFile(ctx, filename)
}
