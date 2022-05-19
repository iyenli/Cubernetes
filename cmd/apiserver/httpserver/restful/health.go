package restful

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHealth(ctx *gin.Context) {
	ctx.String(http.StatusOK, "alive")
}
