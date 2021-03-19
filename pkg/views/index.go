package views

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
	}

	defer db.Close()

	var errSQL error
	var rows *sql.Rows

	rows, errSQL = db.Query("SELECT id, title FROM Section")

	if errSQL != nil {
		httpErr(w, errSQL, 404)
	}

	defer rows.Close()

	sections := []types.Section{}

	if rows != nil {
		for rows.Next() {
			s := types.Section{}
			errScan := rows.Scan(&s.ID, &s.Title)

			if errScan != nil {
				httpErr(w, errScan, 404)
			}

			sections = append(sections, s)
		}
	}

	tmpl, errTmpl := template.ParseFiles("web/templates/index.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
	}

	data := types.IndexData{
		Sections: sections,
	}

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errTmpl, 404)
	} 
}