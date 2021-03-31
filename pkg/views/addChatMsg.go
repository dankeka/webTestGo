package views

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dankeka/webTestGo/types"
	"github.com/gin-gonic/gin"
)

func AddChatMsg(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	req := types.AddChatMsgStruct{
		Csrf: r.FormValue("csrf_token"),
		Text: r.FormValue("textMsg"),
		InterlocutorID: r.FormValue("InterlocutorID"),
		UserID: r.FormValue("UserID"),
	}

	okChan := make(chan bool)

	go func(){
		checkCsrf := CheckCsrf(c, req.Csrf)

		if !checkCsrf && <-okChan {
			okChan <- false
		}
	}()

	go func(){
		sessionUserID := SessionUserId(c)

		if sessionUserID == 0 || fmt.Sprintf("%d", sessionUserID) != req.UserID {
			okChan <- false
		}
	}()

	go func(){
		interlocutorID := c.Param("id")

		if interlocutorID != req.InterlocutorID {
			okChan <- false
		}
	}()

	select {
	case okChan <- false:
		httpErr(w, fmt.Errorf("ERROR"), 404)
		return
	default:
		// pass
	}

	close(okChan)

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, 404)
		return
	}

	defer db.Close()

	time := time.Now().Unix()

	_, errExec := db.Exec(
		"INSERT INTO Message(text, user_id, interlocutor_id, date) VALUES ($1, $2, $3, $4)",
		req.Text,
		req.UserID,
		req.InterlocutorID,
		time,
	)

	if errExec != nil {
		httpErr(w, errExec, 404)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/chat/user/%s", req.InterlocutorID), http.StatusSeeOther)
}