package main

import (
	"log"

	serv "github.com/dankeka/webTestGo"
	"github.com/dankeka/webTestGo/pkg/handler"
)


func main() {
	handlers := new(handler.Handler)

	srv := new(serv.Server)

	if err := srv.Run("8080", handlers.InitRouters()); err != nil {
		log.Fatal("error server: ", err.Error())
	}
}
