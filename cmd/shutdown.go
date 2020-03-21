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

type closer interface {
	Close() error
}
type structForClose struct {
	components []closer
}

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func (closing *structForClose) GracefulShutdown(done chan<- bool) error {
	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Println("Recieved", sig, "signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	case <-time.After(10 * time.Second):
		log.Println("Overslept")
		cancel()
	default:
		cancel()
	}

	// Please, replace srv with your server.
	server := &http.Server{
		Addr: ":8080",
	}

	defer func() {
		for _, component := range closing.components {
			component.Close()
		}

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed:%+v", err)
		}
		cancel()
	}()
	close(done)
	return nil
}