package main

import (
	"log"
)

func main() {

	ConnectedComponents := &structForClose{

		//Put connection variables here
	}

	done := make(chan bool)

	err := ConnectedComponents.GracefulShutdown(done)
	if err != nil {
		log.Fatalf("Service graceful shutdown failed: %v", err)
	}

	<-done
	log.Println("Service successfuly shutdown")

}
