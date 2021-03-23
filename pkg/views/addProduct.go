package views

import (
	"html/template"
	"net/http"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)


func AddProductGet(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	checkLogin := CheckLoginUser(c)
	userId := SessionUserId(c)

	if !checkLogin || userId == 0 {
		http.Error(w, "ERROR", 404)
	}

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
	}

	defer db.Close()

	rows, errQuery := db.Query("SELECT id, title FROM Section")

	if errQuery != nil {
		httpErr(w, errQuery, 404)
	}

	defer rows.Close()

	var sections []types.Section

	for rows.Next() {
		var s types.Section

		errScan := rows.Scan(&s.ID, &s.Title)

		if errScan != nil {
			httpErr(w, errScan, 404)
		}

		sections = append(sections, s)
	}

	csrf, errGenerateCsrf := CsrfGenerate(c)

	if errGenerateCsrf != nil {
		httpErr(w, errGenerateCsrf, 404)
	}

	data := types.AddProductPageStruct{
		Sections: sections,
		IsLogin: checkLogin,
		UserId: userId,
		Csrf: csrf,
	}

	tmpl, errTmpl := template.ParseFiles("web/templates/addProduct.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
	}

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
	}
}