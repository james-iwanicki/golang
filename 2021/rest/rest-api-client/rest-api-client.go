package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

var (
	router *gin.Engine
)

type User struct {
	ID int	`json: "ID"`
	Name string `json: "Name"`
	Age int `json: "Age"`
}

func init() {
	router = gin.Default()
}

func add_user_func(c *gin.Context) {
	var user User
	c.Bind(&user)

	json_data, jerr := json.Marshal(user)
	if jerr != nil {
		c.AbortWithError(http.StatusBadRequest, jerr)
		return
	}

	_, perr := http.Post("http://rest-api-server-service/apis/v1/add_user", "application/json", bytes.NewBuffer(json_data))
	if perr != nil {
		c.AbortWithError(http.StatusBadRequest, perr)
		return
	}
}

func get_user_func(c *gin.Context) {
	Name := c.Query("Name")

	parse, perr := url.Parse("http://rest-api-server-service/apis/v1/get_user")
	if perr != nil {
		c.AbortWithError(http.StatusBadRequest, perr)
		return
	}

	values := url.Values{}
	values.Add("Name", Name)
	parse.RawQuery = values.Encode()

	resp, rerr := http.Get(parse.String())
	if rerr != nil {
		c.AbortWithError(http.StatusBadRequest, rerr)
		return
	}

	defer resp.Body.Close()
	data, derr := ioutil.ReadAll(resp.Body)
	if derr != nil {
		c.AbortWithError(http.StatusBadRequest, derr)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(fmt.Sprintf("%s\n", string(data))))
	
}

func delete_user_func(c *gin.Context) {
	Name := c.Query("Name")
	
	parse, perr := url.Parse("http://rest-api-server-service/apis/v1/delete_user")
	if perr != nil {
		c.AbortWithError(http.StatusBadRequest, perr)
		return
	}

	values := url.Values{}
	values.Add("Name", Name)
	parse.RawQuery = values.Encode()
	
	var jdata []byte

	req, rerr := http.NewRequest(http.MethodDelete, parse.String(), bytes.NewBuffer(jdata))
	//resp, rerr := http.Delete(parse.String())
	if rerr != nil {
		c.AbortWithError(http.StatusBadRequest, rerr)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer resp.Body.Close()
}

func main() {
	v := router.Group("apis/v1")
	{
		v.POST("add_user", add_user_func)
		v.GET("get_user", get_user_func)
		v.DELETE("delete_user", delete_user_func)	
	}
	router.Run(":80")
}
