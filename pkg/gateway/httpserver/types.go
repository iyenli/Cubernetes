package httpserver

import "github.com/gin-gonic/gin"

type Handler struct {
	Method     string
	Path       string
	HandleFunc func(ctx *gin.Context)
}
