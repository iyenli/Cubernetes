package httpserver

import (
	cubeconfig "Cubernetes/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Handler struct {
	Method     string
	Path       string
	HandleFunc func(ctx *gin.Context)
}

func Run() {
	_ = os.MkdirAll(cubeconfig.JobFileDir, 0777)
	_ = os.MkdirAll(cubeconfig.ActionFileDir, 0777)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors())

	router.Static("/workflow", path.Join(cubeconfig.StaticDir, "./workflow"))

	handlerList := append(restfulList, watchList...)

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

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, Cookie, Set-Cookie")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
