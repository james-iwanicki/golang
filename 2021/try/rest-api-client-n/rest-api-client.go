package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"crypto/tls"
	"crypto/x509"
	"bytes"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type User struct {
        ID int  `json: "ID"`
        Name string `json: "Name"`
        Age int `json: "Age"`
}

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func get_user_secure_func(c *gin.Context) {
	Name := c.Query("Name")

	//parse, parse_err := url.Parse("https://rest-api-server-service/apis/v1/get_user")
	parse, parse_err := url.Parse("https://rest-api-server-service.jeevan-namespace.svc.cluster.local/apis/v1/get_user")
	if parse_err != nil {
		c.AbortWithError(http.StatusBadRequest, parse_err)
		return
	}

	values := url.Values{}
	values.Add("Name", Name)
	parse.RawQuery = values.Encode()

	var jdata []byte
	req, req_err := http.NewRequest(http.MethodGet, parse.String(), bytes.NewBuffer(jdata))
	if req_err != nil {
		c.AbortWithError(http.StatusBadRequest, req_err)
		return
	}
	
	rootca, rootca_err := ioutil.ReadFile("/etc/secret-volume/ca.crt")
	if rootca_err != nil {
		c.AbortWithError(http.StatusBadRequest, rootca_err)
		return
	}

	kp, kp_err := tls.LoadX509KeyPair("/etc/secret-volume/tls.crt", "/etc/secret-volume/tls.key")
	if kp_err != nil {
		c.AbortWithError(http.StatusBadRequest, kp_err)
		return
	}
	
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(rootca)
	client := &http.Client {
		Transport: &http.Transport {
			TLSClientConfig: &tls.Config {
				RootCAs: certpool,
				Certificates: []tls.Certificate { kp },
				//InsecureSkipVerify: true,
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

func add_user_secure_func(c *gin.Context) {
	var user User
	c.Bind(&user)
	fmt.Printf("Name:%s Age:%d\n", user.Name, user.Age)

	jdata, jerr := json.Marshal(user)
	if jerr != nil {
		c.AbortWithError(http.StatusBadRequest, jerr)
		return
	}

	parse, parse_err := url.Parse("https://rest-api-server-service.jeevan-namespace.svc.cluster.local/apis/v1/add_user")
	if parse_err != nil {
		c.AbortWithError(http.StatusBadRequest, parse_err)
		return
	}

	req, req_err := http.NewRequest(http.MethodPost, parse.String(), bytes.NewBuffer(jdata))
	if req_err != nil {
		c.AbortWithError(http.StatusBadRequest, req_err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	ca_cert, ca_cert_err := ioutil.ReadFile("/etc/secret-volume/ca.crt")
	if ca_cert_err != nil {
		c.AbortWithError(http.StatusBadRequest, ca_cert_err)
		return
	}
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(ca_cert)



	kp, kp_err := tls.LoadX509KeyPair("/etc/secret-volume/tls.crt", "/etc/secret-volume/tls.key")
	if kp_err != nil {
		c.AbortWithError(http.StatusBadRequest, kp_err)
		return
	}

	client := &http.Client{
		Transport : &http.Transport {
			TLSClientConfig : &tls.Config {
				RootCAs: certpool,
				Certificates: []tls.Certificate { kp },
			},
		},
	}
	resp, resp_err := client.Do(req)
        if resp_err != nil {
                c.AbortWithError(http.StatusBadRequest, resp_err)
                return
        }
        defer resp.Body.Close()
	
}

func delete_user_secure_func(c *gin.Context) {
	Name := c.Query("Name")

	parse, parse_err := url.Parse("https://rest-api-server-service.jeevan-namespace.svc.cluster.local/apis/v1/delete_user")
	if parse_err != nil {
		c.AbortWithError(http.StatusBadRequest, parse_err)
		return
	}
	values := url.Values{}
	values.Add("Name", Name)
	parse.RawQuery = values.Encode()

	var jdata []byte
	req, req_err := http.NewRequest(http.MethodDelete, parse.String(), bytes.NewBuffer(jdata))
	if req_err != nil {
		c.AbortWithError(http.StatusBadRequest, req_err)
		return
	}

	ca_cert, ca_cert_err := ioutil.ReadFile("/etc/secret-volume/ca.crt")
	if ca_cert_err != nil {
		c.AbortWithError(http.StatusBadRequest, ca_cert_err)
		return
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(ca_cert)

	kp, kp_err := tls.LoadX509KeyPair("/etc/secret-volume/tls.crt", "/etc/secret-volume/tls.key")
	if kp_err != nil {
		c.AbortWithError(http.StatusBadRequest, kp_err)
		return
	}

	client := &http.Client {
		Transport: &http.Transport {
			TLSClientConfig: &tls.Config {
				RootCAs: certpool,
				Certificates: []tls.Certificate { kp },
			},
		},
	}

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
		v.GET("/get_user", get_user_secure_func)
		v.POST("/add_user", add_user_secure_func)
		v.DELETE("/delete_user", delete_user_secure_func)
	}
	router.Run(":80")
}
	
