package main

import (
	"log"
	"time"
)

func main() {
	for {
		log.Println(" [INFO] App is running.")
		time.Sleep(5 * time.Second)
	}
}
