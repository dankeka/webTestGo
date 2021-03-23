package handler

import (
	"net/http"

	views "github.com/dankeka/webTestGo/pkg/views"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Handler struct {

}

func (h *Handler) InitRouters() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("auYDg$cxfs2173csjdd423$sdc34su32"))
	r.Use(sessions.Sessions("RunokSessions", store))

	r.GET("/", views.Index)

	register := r.Group("/register")
	{
		register.GET("/get", views.RegisterTmpl)
		register.POST("/post", views.RegisterPost)
	}

	login := r.Group("/login")
	{
		login.GET("/get", views.LoginGet)
		login.POST("/post", views.LoginPost)
	}
	r.GET("/logout", views.Logout)

	user := r.Group("/user")
	{
		user.GET("/me", views.MyUserProfil)
		user.POST("/updateSettings", views.UpdateUserSettings)
	}

	product := r.Group("/product")
	{
		product.GET("/addPage", views.AddProductGet)
	}

	fs := http.Dir("./web/static")

	r.StaticFS("/static/", fs)

	return r
}