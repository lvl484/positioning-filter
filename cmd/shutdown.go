package main

import (
	"context"
	"io"
	"log"
	"time"
)

//GracefulShutdown implements releasing all resouces it got from system, finish all request handling and return responses when service stopping.
func gracefulShutdown(timeout time.Duration, components []io.Closer) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, component := range components {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := component.Close(); err != nil {
				log.Printf("Cant close, err: %v", err)
			}
		}
	}

	return nil
}
