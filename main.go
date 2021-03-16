package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func err(w http.ResponseWriter, err error, cod string) {
	fmt.Println(err.Error())
	w.Write(
		[]byte(fmt.Sprintf("ERROR %s", cod)),
	)
}


func index(w http.ResponseWriter, r *http.Request) {
	tmpl, errTmpl := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if errTmpl != nil {
		err(w, errTmpl, "404")
	}

	tmpl.ExecuteTemplate(w, "index", nil)
}


func body() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", index)

	log.Fatal(http.ListenAndServe(":8080", r))
}


func main() {
	body()
}
