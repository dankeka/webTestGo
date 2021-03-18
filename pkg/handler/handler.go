package handler

import (
	"net/http"

	views "github.com/dankeka/webTestGo/pkg/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {

}

func (h *Handler) InitRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", views.Index)

	fileServer := http.FileServer(http.Dir("./web/static/"))

	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return r
}
