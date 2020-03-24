package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	shutdownTimeout = 10 * time.Second
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func gracefulShutdown(done chan<- bool, components []io.Closer) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	var errTimeout error

	go func() {
		<-ctx.Done()
		errTimeout = ctx.Err()
	}()
	// Please, replace srv with your server.
	server := &http.Server{
		Addr: ":8080",
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}

	for _, component := range components {
		component.Close()
	}
	cancel()

	if errTimeout == context.DeadlineExceeded {
		return errTimeout
	}

	close(done)

	return nil
}
