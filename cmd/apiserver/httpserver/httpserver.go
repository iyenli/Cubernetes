package httpserver

import (
	cubeconfig "Cubernetes/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	Method     string
	Path       string
	HandleFunc func(ctx *gin.Context)
}

func Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	handlerList := append(restfulList, watchList...)
	handlerList = append(handlerList, fileList...)

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
