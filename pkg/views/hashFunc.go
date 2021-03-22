package views

import (
	"crypto/md5"
	"fmt"

	serv "github.com/dankeka/webTestGo"
)


func MD5Encode(text string) string {
	h := md5.New()
	h.Write([]byte(text+serv.SALT))
	return fmt.Sprintf("%x", h.Sum(nil))
}