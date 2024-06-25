package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// cancelOnInterrupt run a thread catching SIGTERM and SIGINT, canceling provided context when they occur.
// After first signal is caught and cancel is called, this thread is closed and
// no more signals are getting called, hence any other interrupt kills the app.
func cancelOnInterrupt(cancel context.CancelFunc) {
	go func() {
		stopSig := make(chan os.Signal, 1)
		signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(stopSig)

		<-stopSig
		cancel()
	}()
}
