package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	serv "github.com/dankeka/webTestGo"
	"github.com/dankeka/webTestGo/pkg/handler"
)


func main() {
	handlers := new(handler.Handler)

	srv := new(serv.Server)

	go func() {
		if err := srv.Run("8080", handlers.InitRouters()); err != nil {
			log.Fatal("error server: ", err.Error())
		}
	}()

	fmt.Println("Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	fmt.Println("Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("error %s", err.Error())
	}
}
