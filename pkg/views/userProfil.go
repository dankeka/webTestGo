package views

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)


func MyUserProfil(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	
	var data types.MyUserProfilStruct

	userId := SessionUserId(c)
	checkLogin := CheckLoginUser(c)

	if userId == 0 || !checkLogin {
		data.Access = false
		data.IsLogin = false
		data.User.ID = 0
		data.User.Name = ""
		data.User.Avatar = ""

		tmpl, errTmpl := template.ParseFiles("web/templates/myUserProfil.html", "web/templates/default.html")

		if errTmpl != nil {
			httpErr(w, errTmpl, 404)
		}
	
		errRenderTmpl := tmpl.Execute(w, data)
	
		if errRenderTmpl != nil {
			httpErr(w, errRenderTmpl, 404)
		}
	} else {
		data.Access = true
		data.IsLogin = true
		data.User.ID = userId
	}

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT name, avatar, email, pub_email, site, age, about_me FROM User WHERE id=$1", userId)
	
	errScan := row.Scan(&data.User.Name, &data.User.Avatar, &data.User.Email, &data.User.PubEmail, &data.User.Site, &data.User.Age, &data.User.AboutMe)

	if errScan != nil {
		httpErr(w, errScan, 404)
		return
	}

	tmpl, errTmpl := template.ParseFiles("web/templates/myUserProfil.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
	}

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
	}
}


func UpdateUserSettings(c *gin.Context) {
	var r *http.Request = c.Request
	var w http.ResponseWriter = c.Writer

	checkUserLogin := CheckLoginUser(c)
	userId := SessionUserId(c)

	if !checkUserLogin || userId == 0 {
		http.Error(w, "error: user is not login", 404)
		return
	}

	form := types.UpdateUserSettingsFormStruct{
		Age: r.FormValue("userAge"),
		Site: r.FormValue("userSite"),
		Email: r.FormValue("userEmail"),
		PubEmail: r.FormValue("pubEmail"),
		AboutMe: r.FormValue("aboueMeTextarea"),
	}

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	defer db.Close()

	var userAge int
	var errAtoi error
	if form.Age != "" {
		userAge, errAtoi = strconv.Atoi(form.Age)

		if errAtoi != nil {
			httpErr(w, errAtoi, 404)
			return
		} 
	}

	var errExec error

	switch {
	case form.Email != "":
		_, errExec = db.Exec(
			"UPDATE User SET email=$1 WHERE id=$2",
			form.Email,
			userId,
		)

		if errExec != nil {
			httpErr(w, errExec, 404)
			return
		}
		fallthrough
	
	case form.PubEmail != "":
		fmt.Println(form.PubEmail)
		if form.PubEmail == "1" {
			_, errExec = db.Exec(
				"UPDATE User SET pub_email=1 WHERE id=$2",
				userId,
			)
		} else if form.PubEmail == "0" {
			_, errExec = db.Exec(
				"UPDATE User SET pub_email=0 WHERE id=$2",
				userId,
			)
		}
	
		if errExec != nil {
			httpErr(w, errExec, 404)
			return
		}
		fallthrough

	case form.Site != "":
		_, errExec = db.Exec(
			"UPDATE User SET site=$1 WHERE id=$2",
			form.Site,
			userId,
		)

		if errExec != nil {
			httpErr(w, errExec, 404)
			return
		}
		fallthrough

	case form.Age != "": 
		_, errExec = db.Exec(
			"UPDATE User SET age=$1 WHERE id=$2",
			userAge,
			userId,
		)

		if errExec != nil {
			httpErr(w, errExec, 404)
			return
		}
		fallthrough

	case form.AboutMe != "":
		_, errExec = db.Exec(
			"UPDATE User SET about_me=$1 WHERE id=$2",
			form.AboutMe,
			userId,
		)

		if errExec != nil {
			httpErr(w, errExec, 404)
			return
		}
	}

	http.Redirect(w, r, "/user/me", http.StatusSeeOther)
}