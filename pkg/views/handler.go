package views

import (
	"net/http"

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

	r.GET("/", h.Index)
	
	register := r.Group("/register")
	{
		register.GET("/get", h.RegisterTmpl)
		register.POST("/post", h.RegisterPost)
	}

	login := r.Group("/login")
	{
		login.GET("/get", h.LoginGet)
		login.POST("/post", h.LoginPost)
	}
	r.GET("/logout", h.Logout)

	user := r.Group("/user")
	{
		user.GET("/me", h.MyUserProfil)
		user.GET("/i/:id", h.OpenUser)
		user.POST("/updateSettings", h.UpdateUserSettings)
		user.POST("/updateAva", h.UpdateAva)
	}

	product := r.Group("/product")
	{
		product.GET("/addPage", h.AddProductGet)
		product.GET("/my", h.MyProducts)
		product.GET("/i/:id", h.OpenProduct)
		product.POST("/addPOST", h.AddProductPost)
	}

	chat := r.Group("/chat")
	{
		chat.GET("/me/chats", h.MyChats)
		chatUser := chat.Group("/user")
		{
			chatUser.GET("/:id", h.OpenChat)
			chatUser.POST("/:id/addMessage", h.AddChatMsg)
		}
	}

	fs := http.Dir("./web/static")

	r.StaticFS("/static/", fs)

	return r
}