package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"net/http"
	"crypto/tls"
	"fmt"
	"crypto/x509"
	"io/ioutil"
)

type User struct {
	ID int `json "ID"`
	Name string `json "Name"`
	Age int `json "Age"`
}

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func get_user_secure_func(c *gin.Context) {
	Name := c.Query("Name")

	db, db_err := sql.Open("mysql", "root:password@tcp(mysql:3306)/test")
	if db_err != nil {
		c.AbortWithError(http.StatusBadRequest, db_err)
		return
	}
	defer db.Close()

	stmt, stmt_err := db.Prepare("select Age from users where Name=?")
	if stmt_err != nil {
		c.AbortWithError(http.StatusBadRequest, stmt_err)
		return
	}

	rows, rows_err  := stmt.Query(Name)
	if rows_err != nil {
		c.AbortWithError(http.StatusBadRequest, rows_err)
		return
	}
	
	var Age int
	for rows.Next() {
		err := rows.Scan(&Age)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(fmt.Sprintf("%s %d\n", Name, Age)))	
	}

}

func add_user_secure_func(c *gin.Context) {
	var user User
	c.Bind(&user)
	fmt.Printf("jejej %s %d\n", user.Name, user.Age)

	db, db_err := sql.Open("mysql", "root:password@tcp(mysql:3306)/test")
	if db_err != nil {
		c.AbortWithError(http.StatusBadRequest, db_err)
		return
	}
	defer db.Close()

	stmt, stmt_err := db.Prepare("insert users set Name=?, Age=?")
	if stmt_err != nil {
		c.AbortWithError(http.StatusBadRequest, stmt_err)
		return
	}
	
	rows, rows_err := stmt.Query(user.Name, user.Age)
	if rows_err != nil {
		c.AbortWithError(http.StatusBadRequest, rows_err)
		return
	}

	err := rows.Err()
        if err != nil {
                c.AbortWithError(http.StatusBadRequest, err)
                return
        }
}

func delete_user_secure_func(c *gin.Context) {
	Name := c.Query("Name")

	db, db_err := sql.Open("mysql", "root:password@tcp(mysql:3306)/test")
        if db_err != nil {
                c.AbortWithError(http.StatusBadRequest, db_err)
                return
        }
        defer db.Close()


	stmt, stmt_err := db.Prepare("delete from users where Name=?")
        if stmt_err != nil {
                c.AbortWithError(http.StatusBadRequest, stmt_err)
                return
        }

	rows, rows_err := stmt.Query(Name)
	if rows_err != nil {
                c.AbortWithError(http.StatusBadRequest, rows_err)
                return
        }

	err := rows.Err()
	if err != nil {
                c.AbortWithError(http.StatusBadRequest, err)
                return
        }
}

func main() {
	v := router.Group("/apis/v1")
	{
		v.GET("/get_user", get_user_secure_func)
		v.POST("/add_user", add_user_secure_func)
		v.DELETE("/delete_user", delete_user_secure_func)
	}
	x509KP, x509KP_err := tls.LoadX509KeyPair("/etc/secret-volume/tls.crt", "/etc/secret-volume/tls.key")
	if x509KP_err != nil {
		return
	}

	//rootca, rootca_err := ioutil.ReadFile("/etc/secret-volume/tls.crt")
	rootca, rootca_err := ioutil.ReadFile("/etc/secret-volume/ca.crt")
	if rootca_err != nil {
		return
	}
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(rootca)
	
	server := &http.Server {
		Addr: ":443",
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate { x509KP },
			ClientCAs: certpool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}
	server.ListenAndServeTLS("", "")
}
