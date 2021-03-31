package views

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)

func OpenChat(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	var data types.ChatStruct

	interlocutorID := c.Param("id")
	userId := SessionUserId(c)

	if userId == 0 || fmt.Sprintf("%d", userId) == interlocutorID {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data.IsLogin = true
	data.UserID = userId
	data.InterlocutorID = interlocutorID

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
	}

	defer db.Close()

	rows, errQuery := db.Query(
		"SELECT id, text, interlocutor_id, user_id, date FROM Message WHERE (user_id=$1 AND interlocutor_id=$2) OR (user_id=$2 AND interlocutor_id=$1)",
		userId,
		interlocutorID,
	)

	if errQuery != nil {
		httpErr(w, errQuery, 404)
		return
	}

	defer rows.Close()

	var messages []types.Message

	for rows.Next() {
		var m types.Message

		errScan := rows.Scan(
			&m.ID, &m.Text, &m.InterlocutorID,
			&m.UserID, &m.Date, 
		)

		if errScan != nil {
			httpErr(w, errScan, 404)
			return
		}

		userNameRow := db.QueryRow(
			"SELECT name FROM User WHERE id=$1",
			m.UserID,
		)

		errScan = userNameRow.Scan(&m.UserName)

		if errScan != nil {
			httpErr(w, errScan, 404)
			return
		}

		messages = append(messages, m)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Date.Unix() > messages[j].Date.Unix()
	})

	data.Messages = messages

	csrfToken, errGenerCsrf := CsrfGenerate(c)

	if errGenerCsrf != nil {
		httpErr(w, errGenerCsrf, 404)
		return
	}

	data.Csrf = csrfToken

	tmpl, errTmpl := template.ParseFiles("web/templates/chat.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		httpErr(w, err, 404)
		return
	}
}