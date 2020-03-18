package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func GracefulShutdown(done <-chan os.Signal, srv *http.Server) {

	<-done
	log.Print("Service stopped")

	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer func() {
		// handle extra services here (like closing database, etc.)
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Service shutdown failed:%+v", err)
	}

	log.Print("Service exited properly")
}
