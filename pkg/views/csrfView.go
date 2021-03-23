package views

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CsrfGenerate(c *gin.Context) (string, error) {
	rand.Seed(time.Now().UnixNano())

	var integer int
	if 10000 > 99999 {
		integer = 10000
	} else {
		integer = rand.Intn(99999-10000) + 10000
	}

	h := sha1.New()
	h.Write([]byte( fmt.Sprintf("%d", integer) ))

	var token string = fmt.Sprintf("%x", h.Sum(nil))

	session := sessions.Default(c)

	session.Set("csrf_token", token)

	errSessionSave := session.Save()
	
	if errSessionSave != nil {
		return "", errSessionSave
	}

	return token, nil
}


func CheckCsrf(c *gin.Context, token string) bool {
	session := sessions.Default(c)

	return token == session.Get("csrf_token")
}