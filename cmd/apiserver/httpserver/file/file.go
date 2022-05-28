package file

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFile(ctx *gin.Context, filename string) {
	ctx.File(filename)
}

func PostFile(ctx *gin.Context, filename string) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BadRequest(ctx)
		return
	}

	err = ctx.SaveUploadedFile(file, filename)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.String(http.StatusOK, "succeeded")
}
