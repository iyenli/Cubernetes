package httpserver

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

var port = 8080

func Run() {
	router := gin.Default()

	for _, handler := range handlerList {
		switch handler.RequestType {
		case HTTP_GET:
			router.GET(handler.Path, handler.HandleFunc)
		case HTTP_POST:
			router.POST(handler.Path, handler.HandleFunc)
		case HTTP_PUT:
			router.PUT(handler.Path, handler.HandleFunc)
		case HTTP_DELETE:
			router.DELETE(handler.Path, handler.HandleFunc)
		}
	}

	err := router.Run(":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal(err, "failure when running api http server")
	}
}
