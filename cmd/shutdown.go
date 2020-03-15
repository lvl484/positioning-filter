package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	timeout = 10 * time.Second
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func GracefulShutdown(server *http.Server) {

	done := make(chan os.Signal, 1)
	signal.Notify(done, signals...)

	<-done
	log.Print("Service stopped")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer func() {
		// handle extra services here (like closing database, etc.)
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Service shutdown failed:%+v", err)
	}

	log.Print("Service exited properly")
}
