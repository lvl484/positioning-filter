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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("Overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

	//Replace srv value with your server
	srv := &http.Server{
		Addr: ":8080",
	}

	defer func() {
		// Handle extra services here (like closing database, etc.)

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed:%+v", err)
		}

		cancel()
	}()

	close(done)

	return nil
}
