package httpserver

import (
	"Cubernetes/pkg/gateway/options"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func handlerTemp(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

func GetGatewayRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	return router
}

func Run(router *gin.Engine) {
	router.GET("/cubernetes/healthTest", handlerTemp)

	err := router.Run(":" + strconv.Itoa(options.GatewayPort))
	if err != nil {
		log.Fatal(err, "[Error]: failure when running Gateway")
	}
}
