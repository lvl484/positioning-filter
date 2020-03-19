package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal)
	done := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err := GracefulShutdown(sigs, done)
	if err != nil {
		log.Fatalf("Service graceful shutdown failed: %v", err)
	}

	<-done
	log.Println("Service successful shutdown")

}
