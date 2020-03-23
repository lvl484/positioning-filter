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

const (
	shutdownTimeout = 10 * time.Second
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func gracefulShutdown(done chan<- bool) error {
	sigs := make(chan os.Signal)

	timeout := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Println("Recieved", sig, "signal")

	ctx, cancel := context.WithCancel(context.Background())

	shutdownTime := time.Now().Add(shutdownTimeout)

	go func() {
		for time.Now().Before(shutdownTime) {
			time.Sleep(time.Second)
			continue
		}
		timeout <- true
	}()

	select {
	case <-timeout:
		log.Println("Overslept")
		cancel()
	default:

		// Please, replace srv with your server.
		server := &http.Server{
			Addr: ":8080",
		}

		for _, component := range components {
			component.Close()
		}

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed:%+v", err)
		}
		cancel()

		close(done)
		close(timeout)
		return nil
	}
	return nil
}
