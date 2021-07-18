package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"crypto/tls"
	"crypto/x509"
)

type User struct {
	ID int `json: "ID"`
	Name string `json: "Name"`
	Age int `json: "Age"`
}

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func get_user_func_secure(c *gin.Context) {
	Name := c.Query("Name")

	parse, parse_err := url.Parse("https://rest-api-server-service/apis/v1/get_user")
	if parse_err != nil {
		c.AbortWithError(http.StatusBadRequest, parse_err)
		return
	}
	values := url.Values{}
	values.Add("Name", Name)
	parse.RawQuery = values.Encode()

	var jbyte []byte
	req, req_err := http.NewRequest(http.MethodGet, parse.String(), bytes.NewBuffer(jbyte))
	if req_err != nil {
                c.AbortWithError(http.StatusBadRequest, req_err)
                return
        }

	cacert, cacert_err := ioutil.ReadFile("/etc/secret-volume/tls.crt")
	if cacert_err != nil {
		c.AbortWithError(http.StatusBadRequest, cacert_err)
		return
	}
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(cacert)

	client := &http.Client {
		Transport: &http.Transport {
			TLSClientConfig: &tls.Config {
				RootCAs: certpool,
				InsecureSkipVerify: true,
			},
		},
	}
	resp, resp_err := client.Do(req)
	if resp_err != nil {
                c.AbortWithError(http.StatusBadRequest, resp_err)
                return
        }
        defer resp.Body.Close()

	data, data_err := ioutil.ReadAll(resp.Body)
	if data_err != nil {
		c.AbortWithError(http.StatusBadRequest, data_err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(fmt.Sprintf("%s\n", string(data))))
}

func get_user_func(c *gin.Context) {
	Name := c.Query("Name")
	
	parse, parse_err := url.Parse("http://rest-api-server-service/apis/v1/get_user")
	if parse_err != nil {
		c.AbortWithError(http.StatusBadRequest, parse_err)
		return
	}

	values := url.Values{}
	values.Add("Name", Name)
	parse.RawQuery = values.Encode()

	resp, resp_err := http.Get(parse.String())
	if resp_err != nil {
		c.AbortWithError(http.StatusBadRequest, resp_err)
		return
	}
	defer resp.Body.Close()

	data, data_err := ioutil.ReadAll(resp.Body)
	if data_err != nil {
		c.AbortWithError(http.StatusBadRequest, data_err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(fmt.Sprintf("%s\n", string(data))))
}

func add_user_func(c *gin.Context) {
	var user User
	c.Bind(&user)

	jdata, jerr := json.Marshal(user)
	if jerr != nil {
		c.AbortWithError(http.StatusBadRequest, jerr)
		return
	}

	resp, resp_err := http.Post("http://rest-api-server-service/apis/v1/add_user", "application/json", bytes.NewBuffer(jdata))
	if resp_err != nil {
		c.AbortWithError(http.StatusBadRequest, resp_err)
		return
	}
	defer resp.Body.Close()

}

func delete_user_func(c *gin.Context) {
	Name := c.Query("Name")
	
	parse, parse_err := url.Parse("http://rest-api-server-service/apis/v1/delete_user")
	if parse_err != nil {
		c.AbortWithError(http.StatusBadRequest, parse_err)
		return
	}

	values := url.Values{}
	values.Add("Name", Name)
	parse.RawQuery = values.Encode()

	var jbyte []byte
	req, req_err := http.NewRequest(http.MethodDelete, parse.String(), bytes.NewBuffer(jbyte))
	if req_err != nil {
                c.AbortWithError(http.StatusBadRequest, req_err)
                return
        }

	client := &http.Client{}
	resp, resp_err := client.Do(req)
	if resp_err != nil {
                c.AbortWithError(http.StatusBadRequest, resp_err)
                return
        }
        defer resp.Body.Close()

}

func main() {
	v := router.Group("/apis/v1")
	{
		v.GET("/get_user", get_user_func_secure)
		v.POST("/add_user", add_user_func)
		v.DELETE("/delete_user", delete_user_func)
	}
	router.Run(":80")
}
