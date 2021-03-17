package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)


func conn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./sait.db")

	//db.SetMaxOpenConns(1)

	if err != nil {
		return db, err
	}

	return db, nil
}


func httpErr(w http.ResponseWriter, err error, cod string) {
	fmt.Println(err.Error())
	w.Write(
		[]byte(fmt.Sprintf("ERROR %s", cod)),
	)
}


func index(w http.ResponseWriter, r *http.Request) {
	db, errConn := conn()

	if errConn != nil {
		httpErr(w, errConn, "404")
	}

	defer db.Close()

	tmpl, errTmpl := template.ParseFiles("templates/index.html", "templates/default.html")

	if errTmpl != nil {
		httpErr(w, errTmpl, "404")
	}

	errRenderTmpl := tmpl.Execute(w, nil)

	if errRenderTmpl != nil {
		httpErr(w, errTmpl, "404")
	} 
}


func body() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", index)

	fileServer := http.FileServer(http.Dir("./static/"))

	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	log.Fatal(http.ListenAndServe(":8080", r))
}


func main() {
	body()
}
