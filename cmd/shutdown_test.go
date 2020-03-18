package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func Endpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Test is what we usually do"))
}

func TestGracefulShutdown(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/test", Endpoint).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	type ShutdownStruct struct {
		signal   <-chan os.Signal
		server   *http.Server
		expected int
	}

	var a <-chan os.Signal = 2
	var ShutdownResults = []ShutdownStruct{
		{a, srv, 1},
	}

	for _, test := range ShutdownResults {
		result := GracefulShutdown(signal, server)
		if result != test.expected {
			t.Fatal("Expected result not given")
		}
	}
}
