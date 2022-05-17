package httpserver

import (
	"Cubernetes/cmd/apiserver/httpserver/utils"
	cubeconfig "Cubernetes/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path"
)

var fileList = []Handler{
	{http.MethodGet, "/apis/gpuJob/file/:uid", getJobFile},
	{http.MethodPost, "/apis/gpuJob/file/:uid", postJobFile},

	{http.MethodGet, "/apis/gpuJob/output/:uid", getJobOutput},
	{http.MethodPost, "/apis/gpuJob/output/:uid", postJobOutput},
}

func getJobFile(ctx *gin.Context) {
	filename := path.Join(cubeconfig.JobFileDir, ctx.Param("uid")+".tar.gz")
	ctx.File(filename)
}

func postJobFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BadRequest(ctx)
		return
	}

	filename := path.Join(cubeconfig.JobFileDir, ctx.Param("uid")+".tar.gz")
	err = ctx.SaveUploadedFile(file, filename)

	if err != nil {
		utils.ServerError(ctx)
		return
	}

	ctx.String(http.StatusOK, "succeeded")
}

func getJobOutput(ctx *gin.Context) {
	filename := path.Join(cubeconfig.JobFileDir, ctx.Param("uid")+".out")
	fmt.Println(filename)
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		ctx.String(http.StatusNotFound, "no job output found")
		return
	}

	ctx.String(http.StatusOK, string(buf))
}

func postJobOutput(ctx *gin.Context) {
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
