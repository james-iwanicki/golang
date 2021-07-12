package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"
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

func get_user_func(c *gin.Context) {
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

	rows, rows_err := stmt.Query(Name)
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
		c.Writer.Write([]byte(fmt.Sprintf("%s %d\n", Name, Age )))
	}
}

func add_user_func(c *gin.Context) {
	var user User
	c.Bind(&user)
	
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

	_, rows_err := stmt.Query(user.Name, user.Age)
	if rows_err != nil {
		c.AbortWithError(http.StatusBadRequest, rows_err)
		return
	}
}

func delete_user_func(c *gin.Context) {
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

	_, rows_err := stmt.Query(Name)
	if rows_err != nil {
		c.AbortWithError(http.StatusBadRequest, rows_err)
		return
	}
}

func main() {
	v := router.Group("/apis/v1")
	{
		v.GET("/get_user", get_user_func)
		v.POST("/add_user", add_user_func)
		v.DELETE("/delete_user", delete_user_func)
	}
	router.Run(":80")
}
