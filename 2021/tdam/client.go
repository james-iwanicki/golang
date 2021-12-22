package main

import (
        "github.com/gin-gonic/gin"
        "net/http"
        "net/url"
	"time"
        "fmt"
        "io/ioutil"
	"bytes"
)

var (
        router *gin.Engine
)

func init() {
        router = gin.Default()
}


func get_user_secure_func(c *gin.Context) {
        //Name := c.Query("Name")

        //parse, parse_err := url.Parse("http://www.google.com")
        parse, parse_err := url.Parse("https://auth.tdameritrade.com/auth")
        //parse, parse_err := url.Parse("https://rest-api-server-service.jeevan-namespace.svc.cluster.local/apis/v1/get_user")
        if parse_err != nil {
                c.AbortWithError(http.StatusBadRequest, parse_err)
                return
        }

        values := url.Values{}
        values.Add("response_type", "code")
        values.Add("redirect_uri", "https://127.0.0.1")
        values.Add("client_id", "JH6HUNVUC1AO9PUIG0EVQDRKKATUNLGI@AMER.OAUTHAP")
        parse.RawQuery = values.Encode()
	fmt.Printf("%s\n", parse.String())

        var jdata []byte
        req, req_err := http.NewRequest(http.MethodGet, parse.String(), bytes.NewBuffer(jdata))
        if req_err != nil {
                c.AbortWithError(http.StatusBadRequest, req_err)
                return
        }
	client := http.Client{Timeout: time.Duration(1) * time.Second}

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

func callback_func(c *gin.Context) {
        c.Writer.WriteHeader(http.StatusOK)
        c.Writer.Write([]byte(fmt.Sprintf("SSP\n")))

}

func main() {
        v := router.Group("/apis/v1")
        {
                v.GET("/get_user", get_user_secure_func)
		v.GET("/callback", callback_func)
        }
        router.Run(":80")
}

