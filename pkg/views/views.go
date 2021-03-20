package views

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)


func conn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "sait.db")

	if err != nil {
		return db, err
	}

	return db, nil
}


func httpErr(w http.ResponseWriter, err error, cod int) {
	fmt.Println(err.Error())
	w.Write(
		[]byte(fmt.Sprintf("ERROR %d", cod)),
	)
}


func CheckLoginUser(c *gin.Context) bool {
	session := sessions.Default(c)
	var isLogin bool

	if session.Get("UserId") != nil {
		isLogin = true
	} else {
		isLogin = false
	}

	return isLogin;
}