package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Test is what we usually do"))
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/test", TestEndpoint).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	done := make(chan os.Signal, 1)

	//signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	signal.Notify(done, syscall.SIGINT)

	go func() {
		GracefulShutdown(done, srv)
	}()
}
