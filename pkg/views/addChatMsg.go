package views

import (
	"fmt"
	"net/http"

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

		if !checkCsrf {
			okChan <- false
			return
		}

		okChan <- true
	}()

	go func(){
		sessionUserID := SessionUserId(c)

		if sessionUserID == 0 || fmt.Sprintf("%d", sessionUserID) != req.UserID {
			okChan <- false
			return
		}

		okChan <- true
	}()

	go func(){
		interlocutorID := c.Param("id")

		if interlocutorID != req.InterlocutorID {
			okChan <- false
			return
		}

		okChan <- true
	}()

	val, ok := <-okChan

	if !val || !ok {
		httpErr(w, fmt.Errorf("ERROR"), 404)
		return
	}

	db, errConn := conn()

	if errConn != nil {
		httpErr(w, fmt.Errorf("ERROR"), 404)
		return
	}

	defer db.Close()
}