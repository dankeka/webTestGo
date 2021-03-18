package main

import (
	"log"

	"github.com/dankeka/webTestGo/pkg/handler"
	serv "github.com/dankeka/webTestGo/server"
)


func main() {
	handlers := new(handler.Handler)

	srv := new(serv.Server)

	if err := srv.Run("8080", handlers.InitRouters()); err != nil {
		log.Fatal("error server: %s", err.Error())
	}
}
