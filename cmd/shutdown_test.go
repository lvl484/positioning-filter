package main

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestGracefulShutdown(t *testing.T) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	type args struct {
		sigs <-chan os.Signal
		done chan<- bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "first",
			args: args{
				sigs: sigs,
				done: done,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GracefulShutdown(tt.args.sigs, tt.args.done); (err != nil) != tt.wantErr {
				t.Errorf("GracefulShutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
