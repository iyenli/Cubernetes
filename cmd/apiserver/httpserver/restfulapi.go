package httpserver

import (
	"github.com/gin-gonic/gin"
)

type HttpRequestType int

const (
	HTTP_GET    HttpRequestType = 0
	HTTP_POST   HttpRequestType = 1
	HTTP_PUT    HttpRequestType = 2
	HTTP_DELETE HttpRequestType = 3
)

type Handler struct {
	RequestType HttpRequestType
	Path        string
	HandleFunc  func(ctx *gin.Context)
}

var handlerList = [...]Handler{
	{HTTP_GET, "/pod", getPod},
}

func getPod(ctx *gin.Context) {

}
