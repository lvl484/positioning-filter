package main

import (
	"net/http"
	"testing"
)

func TestGracefulShutdown(t *testing.T) {
	type ShutdownStruct struct {
		server   *http.Server
		signal   chan<- bool
		expected int
	}

	var ShutdownResults = []ShutdownStruct{
		{srv, true, 2},
	}

	for _, test := range ShutdownResults {
		result := GracefulShutdown(server, signal)
		if result != test.expected {
			t.Fatal("Expected result not given")
		}
	}
}
