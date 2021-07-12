package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

var (
	router *gin.Engine
)

type User struct {
	ID int	`json: "ID"`
	Name string `json: "Name"`
	Age int	`json: "Age"`
}

func init() {
	router = gin.Default()
}

func add_user_func(c *gin.Context) {
	var user User
	c.Bind(&user)

	db, derr := sql.Open("mysql", "root:password@tcp(mysql:3306)/test")
	if derr != nil {
		c.AbortWithError(http.StatusBadRequest, derr)
		return
	}

	stmt, serr := db.Prepare("insert users set Name=?, Age=?")
	if serr != nil {
		c.AbortWithError(http.StatusBadRequest, serr)
		return
	}

	rows, rerr := stmt.Query(user.Name, user.Age)
	if rerr != nil {
		c.AbortWithError(http.StatusBadRequest, rerr)
		return
	}

	err := rows.Err()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	
}

func get_user_func(c *gin.Context) {
	Name := c.Query("Name")

	db, derr := sql.Open("mysql", "root:password@tcp(mysql:3306)/test")
	if derr != nil {
		c.AbortWithError(http.StatusBadRequest, derr)
		return
	}
	defer db.Close()

	stmt, serr := db.Prepare("select age from users where Name=?")
	if serr != nil {
		c.AbortWithError(http.StatusBadRequest, serr)
		return
	}

	rows, rerr := stmt.Query(Name)
	if rerr != nil {
		c.AbortWithError(http.StatusBadRequest, rerr)
		return
	}

	var Age int
	for rows.Next() {
		err := rows.Scan(&Age)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, rerr)
			return
		}
		fmt.Printf("%d\n", Age)
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(fmt.Sprintf("Age:%d\n", Age)))
	}
}

func delete_user_func(c *gin.Context) {
	Name := c.Query("Name")

	db, derr := sql.Open("mysql", "root:password@tcp(mysql:3306)/test")
	if derr != nil {
		c.AbortWithError(http.StatusBadRequest, derr)
		return
	}
	defer db.Close()

	stmt, serr := db.Prepare("delete from users where Name=?")
	if serr != nil {
		c.AbortWithError(http.StatusBadRequest, serr)
		return
	}

	_, rerr := stmt.Query(Name)
	if rerr != nil {
		c.AbortWithError(http.StatusBadRequest, rerr)
		return
	}
	
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
