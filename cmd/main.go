package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shop/pkg/server"
	"syscall"
)

func main() {
	port := ":8080"
	handler, err := server.GetApp()
	if err != nil {
		log.Fatalf("get app error: [%s]", err.Error())
	}
	server := http.Server{
		Addr:    port,
		Handler: handler,
	}

	go func() {
		log.Println("start server on", port)
		server.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	server.Shutdown(context.Background())
	log.Println("server stopped")

}
