package httpserver

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/utils/etcd_helper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

var etcd etcd_helper.ETCDContext

func Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	etcd = etcd_helper.ETCDContext{Client: etcd_helper.NewETCDClient()}
	defer etcd_helper.CloseETCDClient(etcd.Client)

	for _, handler := range handlerList {
		switch handler.Method {
		case http.MethodGet:
			router.GET(handler.Path, handler.HandleFunc)
		case http.MethodPost:
			router.POST(handler.Path, handler.HandleFunc)
		case http.MethodPut:
			router.PUT(handler.Path, handler.HandleFunc)
		case http.MethodDelete:
			router.DELETE(handler.Path, handler.HandleFunc)
		}
	}

	err := router.Run(":" + strconv.Itoa(cubeconfig.APIServerPort))
	if err != nil {
		log.Fatal(err, "failure when running api http server")
	}
}
