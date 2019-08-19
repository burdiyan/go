package mainutil

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// TrapSignals creates a context that gets canceled when SIGTERM or SIGINT is received.
func TrapSignals() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		signal.Stop(ch)
		cancel()
	}()

	return ctx
}

// Run runs fn and checks for error.
func Run(fn func() error) {
	if err := fn(); err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}
