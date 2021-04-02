package views

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)

func (h *Handler) MyChats(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	userId := SessionUserId(c)
	if userId == 0 {
		httpErr(w, fmt.Errorf("ERROR"), 404)
		return
	}

	db, errConn := conn()
	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	defer db.Close()

	rows, errQuery := db.Query(
		"SELECT DISTINCT interlocutor_id, user_id FROM Message WHERE user_id=$1 OR interlocutor_id=$1",
		userId,
	)
	if errQuery != nil {
		httpErr(w, errQuery, 404)
		return
	}

	defer rows.Close()

	var data types.MyChatsStruct
	var chats []types.Chat

	for rows.Next() {
		var interlocutorIdInRow uint
		var userIdInRow uint

		if err := rows.Scan(&interlocutorIdInRow, &userIdInRow); err != nil {
			httpErr(w, err, 404)
			return
		}

		var row *sql.Row
		if interlocutorIdInRow != uint(userId) {
			var con bool
			for _, v := range chats {
				if v.InterlocutorID == interlocutorIdInRow {
					con = true
				}
			}

			if con {
				continue
			}

			row = db.QueryRow(
				"SELECT name FROM User WHERE id=$1",
				interlocutorIdInRow,
			)
		} else {
			interlocutorIdInRow, userIdInRow = userIdInRow, interlocutorIdInRow

			var con bool
			for _, v := range chats {
				if v.InterlocutorID == interlocutorIdInRow {
					con = true
				}
			}

			if con {
				continue
			}

			row = db.QueryRow(
				"SELECT name FROM User WHERE id=$1",
				interlocutorIdInRow,
			)
		}

		var name string

		if err := row.Scan(&name); err != nil {
			httpErr(w, err, 404)
			return
		}

		sqlCode := fmt.Sprintf(
			"%s %s %s",
			"SELECT text FROM Message WHERE",
			"id=(SELECT MAX(id) FROM Message WHERE",
			"(interlocutor_id=$1 AND user_id=$2) OR (interlocutor_id=$2 AND user_id=$1))",
		)

		fmt.Println(sqlCode)

		row = db.QueryRow(
			sqlCode,
			interlocutorIdInRow,
			userIdInRow,
		)

		var lastMsgText sql.NullString

		if err := row.Scan(&lastMsgText); err != nil {
			httpErr(w, err, 404)
			return
		}

		if lastMsgText.String == "" || !lastMsgText.Valid {
			continue
		}

		c := types.Chat{
			InterlocutorID: interlocutorIdInRow,
			ChatUserName: name,
			LastMsg: lastMsgText.String,
		}

		chats = append(chats, c)
	}
	
	data.Chats = chats
	data.IsLogin = true

	tmpl, errTmpl := template.ParseFiles("web/templates/myChats.html", "web/templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, 404)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		httpErr(w, err, 404)
	}
}