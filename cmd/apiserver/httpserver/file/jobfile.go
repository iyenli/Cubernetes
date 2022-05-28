package file

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	cubeconfig "Cubernetes/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path"
)

func GetJobFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.JobFileDir, ctx.Param("uid")+".tar.gz")
	ctx.File(filename)
}

func PostJobFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.JobFileDir, ctx.Param("uid")+".tar.gz")
	PostFile(ctx, filename)
}

func GetJobOutput(ctx *gin.Context) {
	filename := path.Join(cubeconfig.JobFileDir, ctx.Param("uid")+".out")
	fmt.Println(filename)
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		ctx.String(http.StatusNotFound, "no job output found")
		return
	}

	ctx.String(http.StatusOK, string(buf))
}

func PostJobOutput(ctx *gin.Context) {
	output, err := ctx.GetRawData()
	if err != nil {
		utils.BadRequest(ctx)
		return
	}

	filename := path.Join(cubeconfig.JobFileDir, ctx.Param("uid")+".out")
	err = ioutil.WriteFile(filename, output, 0600)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.String(http.StatusOK, "succeeded")
}
