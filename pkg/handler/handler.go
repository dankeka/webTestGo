package handler

import (
	"net/http"

	views "github.com/dankeka/webTestGo/pkg/views"
	"github.com/gin-gonic/gin"
)

type Handler struct {

}

func (h *Handler) InitRouters() *gin.Engine {
	r := gin.Default()

	r.GET("/", views.Index)

	register := r.Group("/register")
	{
		register.GET("/get", views.RegisterTmpl)
		//register.POST("/post", nil)
	}

	fs := http.Dir("./web/static")

	r.StaticFS("/static/", fs)

	return r
}