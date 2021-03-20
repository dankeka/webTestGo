package views

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	serv "github.com/dankeka/webTestGo"
	"github.com/dankeka/webTestGo/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetErrLoginCookie(c *gin.Context, value string) {
	cookie, err := c.Cookie("loginErr")
	
	if cookie == "" || err != nil {
		c.SetCookie("loginErr", value, 3600, "/", serv.DOMAIN, false, false)
	}
}


func DeleteErrLoginCookie(c *gin.Context) {
	c.SetCookie("loginErr", "", -1, "/", serv.DOMAIN, false, false)
}


// GET
func LoginGet(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	tmpl, errTmpl := template.ParseFiles("web/templates/login.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
	}

	data := types.LoginGetStruct{}

	checkRegisterErr, errCookie := c.Cookie("loginErr")
	if checkRegisterErr == "" || errCookie != nil {
		data.ErrLogin = false
		data.ErrText = ""
		DeleteErrLoginCookie(c)
	} else {
		data.ErrLogin = true
		data.ErrText = checkRegisterErr
		DeleteErrLoginCookie(c)
	}
	data.IsLogin = CheckLoginUser(c)

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
	}
}


// POST
func LoginPost(c *gin.Context) {
	var r *http.Request = c.Request
	var w http.ResponseWriter = c.Writer

	user := types.LoginPostStruct{
		NameOrEmail: r.FormValue("nameOrEmail"),
		Password: r.FormValue("password"),
	}

	isEmail, errRegex := regexp.MatchString(`.+@.+\..+`, user.NameOrEmail)

	if errRegex != nil {
		fmt.Println(errRegex.Error())
		SetErrLoginCookie(c, "Произошла ошибка! Попробуйте ввести данные ещё раз или попробуйте позже!")
		http.Redirect(w, r, "/login/get", http.StatusSeeOther)
		return
	}

	db, errConn := conn()

	if errConn != nil {
		SetErrLoginCookie(c, "Произошла ошибка! Попробуйте позже!")
		http.Redirect(w, r, "/login/get", http.StatusSeeOther)
		return
	}

	defer db.Close()

	var row *sql.Row
	var errScan error
	var userId sql.NullInt32

	if isEmail {
		row = db.QueryRow("SELECT id FROM User WHERE email=$1 AND password=$2", user.NameOrEmail, user.Password)
		errScan = row.Scan(&userId)
	} else {
		row = db.QueryRow("SELECT id FROM User WHERE name=$1 AND password=$2", user.NameOrEmail, user.Password)
		errScan = row.Scan(&userId)
	}

	if errScan != nil || errScan == sql.ErrNoRows || userId.Int32 == 0 {
		SetErrLoginCookie(c, "Произошла ошибка! Попробуйте ввести данные ещё раз или попробуйте позже!")
		http.Redirect(w, r, "/login/get", http.StatusSeeOther)
		return
	}

	session := sessions.Default(c)
	session.Set("UserId", int(userId.Int32))
	errSessionSave := session.Save()

	if errSessionSave != nil {
		SetErrLoginCookie(c, "Произошла ошибка! Попробуйте позже!")
		http.Redirect(w, r, "/login/get", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GET
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("UserId")
	session.Save()

	var r *http.Request = c.Request
	var w http.ResponseWriter = c.Writer

	http.Redirect(w, r, "/", http.StatusSeeOther)
}