package main

import (
	"io"
	"testing"
)

type CloserMock struct {
	name     string
	expected func() error
}

func NewCloserMock(name string, expected func() error) io.Closer {
	return CloserMock{
		name:     name,
		expected: expected,
	}
}

func (cm CloserMock) Close() error {
	return cm.expected()
}

func TestGracefulShutdownSuccess(t *testing.T) {
	closers := []io.Closer{
		NewCloserMock("consul", func() error {
			return nil
		}),
		NewCloserMock("postgres", func() error {
			return nil
		}),
		NewCloserMock("kafka", func() error {
			return nil
		}),
	}

	err := gracefulShutdown(shutdownTimeout, closers)
	if err != nil {
		t.Error(err)
	}
}
