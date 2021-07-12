package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func get_user_func(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(fmt.Sprintf("jeevan\n")))
}


func main() {
	v := router.Group("/apis/v1")
	{
		v.GET("/get_user", get_user_func)
	}
	router.Run(":80")
}
