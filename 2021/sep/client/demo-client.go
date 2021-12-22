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

type User struct {
        ID int  `json: "ID"`
        Name string `json: "Name"`
        Age int `json: "Age"`
}


func init() {
        router = gin.Default()
}

func envoy_calling_func(c *gin.Context) {
        c.Writer.WriteHeader(http.StatusOK)
        c.Writer.Write([]byte(fmt.Sprintf("Envoy called  Date %s\n", time.Now())))

}

func envoy_calling_post_func(c *gin.Context) {
	var user User
        c.Bind(&user)
        c.Writer.WriteHeader(http.StatusOK)
        c.Writer.Write([]byte(fmt.Sprintf("Envoy called  Date %s %s %d\n", time.Now(), user.Name, user.Age)))
	fmt.Printf("Envoy called  Date %s %s %d\n", time.Now(), user.Name, user.Age)
	
}

func main() {
        v := router.Group("/apis/v1")
        {
                v.GET("/envoy_calling", envoy_calling_func)
		v.POST("/envoy_calling_post", envoy_calling_post_func)
        }
        router.Run(":80")
}

