package main

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestGracefulShutdown_Success1(t *testing.T) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err := GracefulShutdown(sigs, done)
	if err != nil {
		t.Error(err)
	}
}

func TestGracefulShutdown_Fail(t *testing.T) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err := GracefulShutdown(sigs, done)
	if err == nil {
		t.Error("want error got nil")
	}
}
