package views

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	serv "github.com/dankeka/webTestGo"
	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)


func SetErrRegisterCookie(c *gin.Context, value string) {
	cookie, err := c.Cookie("registerErr")
	
	if cookie == "" || err != nil {
		c.SetCookie("registerErr", value, 3600, "/", serv.DOMAIN, false, false)
	}
}


func DeleteErrRegisterCookie(c *gin.Context) {
	c.SetCookie("registerErr", "", -1, "/", serv.DOMAIN, false, false)
}


// GET
func RegisterTmpl(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	tmpl, errTmpl := template.ParseFiles("web/templates/register.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
	}

	data := types.RegisterStruct{}

	isLogin := CheckLoginUser(c)

	checkRegisterErr, errCookie := c.Cookie("registerErr")
	if checkRegisterErr == "" || errCookie != nil {
		data.ErrRegister = false
		data.ErrText = ""
		DeleteErrRegisterCookie(c)
	} else {
		data.ErrRegister = true
		data.ErrText = checkRegisterErr
		DeleteErrRegisterCookie(c)
	}
	data.IsLogin = isLogin

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
	}
}


// POST
func RegisterPost(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	username := r.FormValue("name")

	if username == "" {
		fmt.Println("error username")
		SetErrRegisterCookie(c, "Поле имени пустое")
		http.Redirect(w, r, "/register/get", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	if password1 != password2 {
		fmt.Println("error password")
		SetErrRegisterCookie(c, "Пароли не совпадают!")
		http.Redirect(w, r, "/register/get", http.StatusSeeOther)
		return
	}

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT id FROM User WHERE name=$1 OR email=$2", username, email)

	var checkUserName sql.NullInt32
	errScan := row.Scan(&checkUserName)

	if errScan != nil && errScan != sql.ErrNoRows {
		httpErr(w, errScan, 404)
		return
	}

	if checkUserName.Valid {
		SetErrRegisterCookie(c, "Пользователь с таким именем или эл. почтой уже существует!")
		http.Redirect(w, r, "/register/get", http.StatusSeeOther)
		return
	}

	hashPassword := MD5Encode(password1)

	_, errExec := db.Exec("INSERT INTO User(name, email, password) VALUES ($1, $2, $3)", username, email, hashPassword)

	if errExec != nil {
		httpErr(w, errExec, 404)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}