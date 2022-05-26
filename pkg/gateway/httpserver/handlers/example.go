package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Hello(ctx *gin.Context) {
	fmt.Println("Hello!")
}
