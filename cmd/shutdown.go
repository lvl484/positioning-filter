package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func GracefulShutdown(server *http.Server, done chan<- bool) {

	log.Print("Service stopped") //тест на сигнал

	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer func() {
		// handle extra services here (like closing database, etc.)
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Service shutdown failed:%+v", err) //тест
	}

	log.Print("Service exited properly")
}
