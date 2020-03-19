package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	for {
		log.Println(" [INFO] App is running.")
		time.Sleep(5 * time.Second)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := GracefulShutdown(sigs, done)
		if err != nil {
			log.Fatalf("Service graceful shutdown failed: %v", err)
		}
	}()

	<-done
	log.Println("Service successful shutdown")
}
