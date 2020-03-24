package main

import (
	"context"
	"io"
	"time"
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func gracefulShutdown(timeout time.Duration, components []io.Closer) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	for _, component := range components {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			cancel()
			return err
		default:
			component.Close()
			if err := component.Close(); err != nil {
				return err
			}
		}
	}

	cancel()

	return nil
}
