package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func GracefulShutdown(sigs <-chan os.Signal, done chan<- bool) error {
	sig := <-sigs
	fmt.Println("Recieved", sig, "signal")

	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	//Replace srv value with your server
	srv := &http.Server{
		Addr: ":8080",
	}

	defer func() {
		// Handle extra services here (like closing database, etc.)
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Service shutdown failed:%+v", err)
	}

	done <- true

	return nil
}
