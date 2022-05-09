package restful

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func parseFail(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, "fail to parse body")
}

func badRequest(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, "bad request")
}

func serverError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, "server error")
}

func notFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, "no objects found")
}
