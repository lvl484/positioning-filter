package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {

	done := make(chan os.Signal, 1)

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	signal.Notify(done, signals...)

	//Please, use GracefulShotdown function to terminate all the running processes and make sure resources are released properly.
	//GracefulShutdown(server,done)
	<-done

}
