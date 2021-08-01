package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"time"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func server_func(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(fmt.Sprintf("Server %s\n", time.Now())))
}

func main() {
	router.GET("/server", server_func)
	router.Run(":80")
}
