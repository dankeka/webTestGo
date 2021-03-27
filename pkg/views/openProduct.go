package views

import (
	"html/template"
	"net/http"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)


func OpenProduct(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	var data types.OpenProductStruct

	isLogin := CheckLoginUser(c)

	if !isLogin {
		http.Error(w, "ERROR", 404)
	}

	data.IsLogin = isLogin

	productID := c.Param("id")

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	defer db.Close()

	row := db.QueryRow(
		"SELECT id, title, description, active, section_id, user_id, date, price FROM Product WHERE id=$1",
		productID,
	)

	errScan := row.Scan(
		&data.Product.ID, &data.Product.Title, 
		&data.Product.Description, &data.Product.Active, 
		&data.Product.SectionID, &data.Product.UserID, 
		&data.Product.Date, &data.Product.Price,
	)

	if errScan != nil {
		httpErr(w, errScan, 404)
		return
	}

	rows, errQuery := db.Query(
		"SELECT src FROM ImageProduct WHERE product_id=$1",
		data.Product.ID,
	)

	if errQuery != nil {
		httpErr(w, errQuery, 404)
	}

	defer rows.Close()

	for rows.Next() {
		var imgUrl string

		errScan = rows.Scan(&imgUrl)

		if errScan != nil {
			httpErr(w, errScan, 404)
		}

		data.ImgUrls = append(data.ImgUrls, imgUrl)
	}

	sellerRow := db.QueryRow(
		"SELECT name, avatar FROM User WHERE id=$1",
		data.Product.UserID,
	)

	var sellerData types.UserAvaAndName

	errScan = sellerRow.Scan(&sellerData.Name, &sellerData.AvaUrl)

	if errScan != nil {
		httpErr(w, errScan, 404)
		return
	}

	data.SellerInfo = sellerData

	tmpl, errTmpl := template.ParseFiles("web/templates/product.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
		return
	}

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
	}
}