package httpserver

import (
	"Cubernetes/pkg/gateway/options"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func handlerTemp(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello")
}

func GetGatewayRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	return router
}

func Run(router *gin.Engine) {

	fmt.Println("Add /hi slowly...")
	go func() {
		for i := 3; i < 100; i++ {
			router.GET(strings.Repeat("/hi", i), handlerTemp)
		}
		time.Sleep(2 * time.Second)
	}()

	router.GET("/hi", handlerTemp)
	err := router.Run(":" + strconv.Itoa(options.GatewayPort))
	if err != nil {
		log.Fatal(err, "[Error]: failure when running Gateway")
	}
}
