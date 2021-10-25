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

func envoy_calling_func(c *gin.Context) {
        c.Writer.WriteHeader(http.StatusOK)
        c.Writer.Write([]byte(fmt.Sprintf("Envoy called  Date %s\n", time.Now())))

}


func main() {
        v := router.Group("/apis/v1")
        {
                v.GET("/envoy_calling", envoy_calling_func)
        }
        router.Run(":80")
}

