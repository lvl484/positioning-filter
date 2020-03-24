package main

import (
	"context"
	"io"
	"time"
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func gracefulShutdown(timeout time.Duration, done chan<- bool, components []io.Closer) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	for _, component := range components {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			cancel()
			return err
		default:
			component.Close()
		}
	}

	cancel()

	close(done)

	return nil
}
