package handler

import (
	views "github.com/dankeka/webTestGo/pkg/views"
	"github.com/gin-gonic/gin"
)

type Handler struct {

}

// func (h *Handler) InitRouters() *chi.Mux {
// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger)

// 	r.Get("/", views.Index)

// 	fileServer := http.FileServer(http.Dir("./web/static/"))

// 	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

// 	return r
// }

func (h *Handler) InitRouters() *gin.Engine {
	r := gin.Default()

	r.GET("/", views.Index)

	r.Static("./web/static/", "/static/")

	return r
}