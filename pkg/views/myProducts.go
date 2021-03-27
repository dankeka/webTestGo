package views

import (
	"html/template"
	"net/http"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)

func MyProducts(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	userId := SessionUserId(c)
	isLogin := CheckLoginUser(c)

	if !isLogin || userId == 0 {
		http.Error(w, "ERROR", 404)
	}

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
	}

	defer db.Close()

	rows, errQuery := db.Query(
		"SELECT id, title, description, active, section_id, user_id, date, price FROM Product WHERE user_id=$1",
		userId,
	)

	if errQuery != nil {
		httpErr(w, errQuery, 404)
	}

	defer rows.Close()

	var products []types.ProductAndImg

	for rows.Next() {
		var p types.ProductAndImg

		errScan := rows.Scan(
			&p.Product.ID, &p.Product.Title, &p.Product.Description, 
			&p.Product.Active, &p.Product.SectionID, &p.Product.UserID, 
			&p.Product.Date, &p.Product.Price,
		)
		
		p.Product.DateStr = p.Product.Date.Format("2006-01-02 15:04")

		if errScan != nil {
			httpErr(w, errScan, 404)
		}

		row := db.QueryRow(
			"SELECT src FROM ImageProduct WHERE id=(SELECT MIN(id) FROM ImageProduct WHERE product_id=$1)", 
			p.Product.ID,
		)

		errScan = row.Scan(&p.ImgUrl)

		if errScan != nil {
			httpErr(w, errScan, 404)
		}

		products = append(products, p)
	}

	data := types.MyProductsStruct{
		IsLogin: isLogin,
		Products: products,
	}

	tmpl, errTmpl := template.ParseFiles("web/templates/myProducts.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
	}

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
	}
}