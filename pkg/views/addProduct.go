package views

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)


func DeleteImages(listName []string) error {
	for name := range listName {
		err := os.Remove( fmt.Sprintf("web/static/images/%s", listName[name]) )

		if err != nil {
			return err
		}
	}

	return nil
}


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


func AddProductPost(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	csrfToken := r.FormValue("csrfToken")

	chechCsrf := CheckCsrf(c, csrfToken)

	if !chechCsrf {
		http.Error(w, "ERROR", 404)
	}

	form, errForm := c.MultipartForm()

	if errForm != nil {
		httpErr(w, errForm, 404)
		return
	}

	formData := types.AddProductPost{
		Title: r.FormValue("title"),
		SectionID: r.FormValue("section"),
		Description: r.FormValue("description"),
		Files: form.File["addImg"],
		Price: r.FormValue("priceProduct"),
	}

	intPrice, errAtoi := strconv.Atoi(formData.Price)

	if errAtoi != nil {
		httpErr(w, errAtoi, 404)
	}

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	var fileNameList []string

	for _, file := range formData.Files {
		row := db.QueryRow("SELECT MAX(id) FROM ImageProduct")
		var maxIdImg sql.NullInt32

		errScan := row.Scan(&maxIdImg)

		if errScan != nil {
			DeleteImages(fileNameList)
			httpErr(w, errScan, 404)
			return
		}

		if maxIdImg.Int32 == 0 {
			maxIdImg.Int32 = 1
		} else {
			maxIdImg.Int32 += 1
		}

		fileName := fmt.Sprintf("product_img_%d", maxIdImg.Int32)

		filePath, _ := filepath.Abs("web/static/images/" + fileName)
		f, errCreateFile := os.Create(filePath)

		if errCreateFile != nil {
			DeleteImages(fileNameList)
			httpErr(w, errCreateFile, 404)
		}

		defer f.Close()

		fo, errOpenFile := file.Open()

		if errOpenFile != nil {
			DeleteImages(fileNameList)
			httpErr(w, errOpenFile, 404)
		}

		fileContent, errReadFile := ioutil.ReadAll(fo)

		if errReadFile != nil {
			DeleteImages(fileNameList)
			httpErr(w, errReadFile, 404)
		}

		f.Write(fileContent)

		fileNameList = append(fileNameList, fileName)
	}

	userId := SessionUserId(c)

	if userId == 0 {
		DeleteImages(fileNameList)
		http.Error(w, "ERROR", 404)
	}

	dateNow := time.Now().UnixNano()

	_, errExec := db.Exec(
		"INSERT INTO Product(title, description, price, section_id, user_id, active, date) VALUES ($1, $2, $3, $4, $5, 1, $6)",
		formData.Title,
		formData.Description,
		intPrice,
		formData.SectionID,
		userId,
		dateNow,
	)

	if errExec != nil {
		DeleteImages(fileNameList)
		httpErr(w, errExec, 404)
	}
}