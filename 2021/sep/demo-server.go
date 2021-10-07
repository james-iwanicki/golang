package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func server_func(c *gin.Context) {
	version := os.Getenv("version")
	node := os.Getenv("MY_NODE_NAME")
	pod := os.Getenv("MY_POD_NAME")
	ip := os.Getenv("MY_POD_IP")

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(fmt.Sprintf("Version %s IP %s Pod %s Node %s Date %s\n", version, ip, pod, node, time.Now())))
}

func head_func(c *gin.Context) {
	c.Next()
}

func main() {
	v := router.Group("/apis/v1") 
	{		
		v.GET("/server", server_func)
		v.HEAD("/server", head_func)
	}
	router.Run(":80")
}
