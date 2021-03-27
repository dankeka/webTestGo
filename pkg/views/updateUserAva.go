package views

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)


func UpdateAva(c *gin.Context) {
	var r *http.Request = c.Request
	var w http.ResponseWriter = c.Writer

	csrfToken := r.FormValue("csrf_token")
	checkCsrf := CheckCsrf(c, csrfToken)

	if !checkCsrf {
		http.Error(w, "ERROR", 404)
	}

	userId :=	SessionUserId(c)

	if userId == 0 {
		http.Error(w, "ERROR", 404)
		return
	}

	_, newAva, errForm := r.FormFile("newAva")

	if errForm != nil {
		httpErr(w, errForm, 404)
		return
	}

	fileName := fmt.Sprintf("user_avatar_%d", userId)

	filePath := "web/static/avatars/" + fileName

	fileAva, errOpen := os.Open(filePath)

	if errOpen != nil {
		var errCreate error

		fileAva, errCreate = os.Create(filePath)

		if errCreate != nil {
			httpErr(w, errCreate, 404)
			return
		}
	}

	defer fileAva.Close()

	newAvaFile, errOpenAva := newAva.Open()

	if errOpenAva != nil {
		httpErr(w, errOpenAva, 404)
		return
	}

	fileContent, errReadFile := ioutil.ReadAll(newAvaFile)

	if errReadFile != nil {
		httpErr(w, errReadFile, 404)
		return
	}

	fileAva.Write(fileContent)

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
	}

	_, errExec := db.Exec(
		"UPDATE User SET avatar=$1 WHERE id=$2",
		fileName,
		userId,
	)

	if errExec != nil {
		httpErr(w, errExec, 404)
	}

	http.Redirect(w, r, "/user/me", http.StatusSeeOther)
}