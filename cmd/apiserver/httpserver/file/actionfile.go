package file

import (
	cubeconfig "Cubernetes/config"
	"github.com/gin-gonic/gin"
	"path"
)

func GetActionFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.ActionFileDir, ctx.Param("name")+".py")
	GetFile(ctx, filename)
}

func PostActionFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.ActionFileDir, ctx.Param("name")+".py")
	PostFile(ctx, filename)
}
