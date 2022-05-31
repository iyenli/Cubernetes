package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()

		router.GET("/", handler)
		err := router.Run(":" + strconv.Itoa(8080))
		if err != nil {
			log.Fatal(err, "failure when running api http server")
		}
	}()

	time.Sleep(80 * time.Second)
	log.Println("Oops...I'll go...")
	os.Exit(1)
}

func handler(ctx *gin.Context) {
	ctx.String(200, "Hello, Cubernetes user, I'll kill myself in 80 seconds...\n")
}
