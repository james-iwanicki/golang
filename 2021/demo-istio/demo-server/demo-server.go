package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"time"
	"os"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func server_func(c *gin.Context) {
	version := os.Getenv("version")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(fmt.Sprintf("Server %s Version %s\n", time.Now(), version)))
}

func main() {
	router.GET("/server", server_func)
	router.Run(":80")
}
