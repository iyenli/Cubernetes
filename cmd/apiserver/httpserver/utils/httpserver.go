package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ParseFail(ctx *gin.Context) {
	ctx.String(http.StatusBadRequest, "fail to parse body")
}

func BadRequest(ctx *gin.Context) {
	ctx.String(http.StatusBadRequest, "bad request")
}

func ServerError(ctx *gin.Context) {
	ctx.String(http.StatusInternalServerError, "server error")
}

func NotFound(ctx *gin.Context) {
	ctx.String(http.StatusNotFound, "no objects found")
}
