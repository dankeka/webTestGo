package views

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	var data types.IndexData

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	defer db.Close()

	var errSQL error
	var rows *sql.Rows

	rows, errSQL = db.Query("SELECT id, title FROM Section")

	if errSQL != nil {
		httpErr(w, errSQL, 404)
		return
	}

	defer rows.Close()

	sections := []types.Section{}

	if rows != nil {
		for rows.Next() {
			s := types.Section{}
			errScan := rows.Scan(&s.ID, &s.Title)

			if errScan != nil {
				httpErr(w, errScan, 404)
				return
			}

			sections = append(sections, s)
		}
	}

	rowsProducts, errQuery := db.Query(
		"SELECT id, title, price FROM Product WHERE id>(SELECT MAX(id) FROM Product)-20",
	)

	if errQuery != nil {
		httpErr(w, errQuery, 404)
		return
	}

	defer rowsProducts.Close()

	var newProducts []types.ProductIdAndTitleAndImg

	for rowsProducts.Next() {
		var i types.ProductIdAndTitleAndImg

		errScan := rowsProducts.Scan(&i.ID, &i.Title, &i.Price)

		if errScan != nil {
			httpErr(w, errScan, 404)
			return
		}

		row := db.QueryRow(
			"SELECT src FROM ImageProduct WHERE id=(SELECT MIN(id) FROM ImageProduct WHERE product_id=$1)",
			i.ID,
		)

		var imgUrl string

		errScan = row.Scan(&imgUrl)

		if errScan != nil {
			httpErr(w, errScan, 405)
			return
		}

		i.ImgUrl = imgUrl

		newProducts = append(newProducts, i)
	}

	isLogin := CheckLoginUser(c)

	tmpl, errTmpl := template.ParseFiles("web/templates/index.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
		return
	}

	data = types.IndexData{
		Sections: sections,
		NewProducts: newProducts,
		IsLogin: isLogin,
	}

	fmt.Println(data.NewProducts)

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
		return
	} 
}
