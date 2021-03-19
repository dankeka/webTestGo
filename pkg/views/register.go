package views

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"reflect"

	//"github.com/snowplow/referer-parser/go"
	"github.com/gin-gonic/gin"
)

// GET
func RegisterTmpl(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	tmpl, errTmpl := template.ParseFiles("web/templates/register.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
	}

	errRenderTmpl := tmpl.Execute(w, nil)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
	}
}

// GET
func getUserIdByNameOrEmail(c *gin.Context) {
	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")

	db, errConn := conn()

	defer db.Close()

	if errConn != nil {
		jsonData := []byte(
			fmt.Sprintf(`{
				"check": %t,
				"error": %t,
			}`, false, true),
		)
		c.Data(http.StatusOK, "application/json", jsonData)

		return
	}

	checkUserRow := db.QueryRow("SELECT id FROM User WHERE name=$1 OR email=$2", name, email)

	var userId sql.NullInt32

	errScan := checkUserRow.Scan(&userId)

	if errScan != nil {
		jsonData := []byte(
			fmt.Sprintf(`{
				"check": %t,
				"error": %t,
			}`, false, true),
		)
		c.Data(http.StatusOK, "application/json", jsonData)

		return
	}

	if reflect.TypeOf(userId) != nil {
		jsonData := []byte(
			fmt.Sprintf(`{
				"check": %t,
				"error": %t,
			}`, false, true),
		)
		c.Data(http.StatusOK, "application/json", jsonData)

		return
	}

	jsonData := []byte(
		fmt.Sprintf(`{
			"check": %t,
			"error": %t,
		}`, true, false),
	)
	c.Data(http.StatusOK, "application/json", jsonData)
}

// POST
func RegisterPost(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	username := r.FormValue("name")
	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	if password1 != password2 {
		http.Redirect(w, r, "/register/get", http.StatusSeeOther)
		return
	}

	db, errConn := conn()

	defer db.Close()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	row := db.QueryRow("SELECT id FROM User WHERE name=$1 OR email=$2", username, email)

	var checkUserName sql.NullInt32
	errScan := row.Scan(&checkUserName)

	if errScan != nil {
		httpErr(w, errScan, 404)
		return
	}

	if reflect.TypeOf(checkUserName) != nil {
		http.Redirect(w, r, "/register/get", http.StatusSeeOther)
		return
	}

	_, errExec := db.Exec("INSERT INTO User(name, email, password) VALUES ($1, $2, $3)", username, email, password1)

	if errExec != nil {
		httpErr(w, errExec, 404)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}