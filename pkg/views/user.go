package views

import (
	"html/template"
	"net/http"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)

func OpenUser(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	userId := c.Param("id")
	myUserId := SessionUserId(c)

	if string(rune(myUserId)) == userId {
		http.Redirect(w, c.Request, "/user/me", http.StatusSeeOther)
		return
	}

	var data types.OpenUserAccStruct

	data.IsLogin = CheckLoginUser(c)

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	defer db.Close()

	row := db.QueryRow(
		"SELECT id, name, avatar, email, pub_email, cite, age, about_me FROM User WHERE id=$1",
		userId,
	)

	var user types.User

	errScan := row.Scan(
		&user.ID, &user.Name, &user.Avatar,
		&user.Email, &user.PubEmail, &user.Cite, 
		&user.Age, &user.AboutMe,
	)

	if errScan != nil {
		httpErr(w, errScan, 404)
		return
	}

	data.User = user

	rows, errQuery := db.Query(
		"SELECT id, title, price FROM Product WHERE user_id=$1",
		userId,
	)

	if errQuery != nil {
		httpErr(w, errQuery, 404)
		return
	}

	defer rows.Close()

	var products []types.ProductIdAndTitleAndImg

	for rows.Next() {
		var p types.ProductIdAndTitleAndImg

		if errScan = rows.Scan(&p.ID, &p.Title, &p.Price); errScan != nil {
			httpErr(w, errScan, 404)
			return
		}

		row = db.QueryRow(
			"SELECT src FROM ImageProduct WHERE id=(SELECT MIN(id) FROM ImageProduct WHERE product_id=$1)",
			p.ID,
		)

		if errScan = row.Scan(&p.ImgUrl); errScan != nil {
			httpErr(w, errScan, 404)
			return
		}

		products = append(products, p)
	}

	data.LastProducts = products

	tmpl, errTmpl := template.ParseFiles("web/templates/user.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
		return
	}

	errRenderTmpl := tmpl.Execute(w, data)

	if errRenderTmpl != nil {
		httpErr(w, errRenderTmpl, 404)
		return
	}
}