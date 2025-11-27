package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to catch OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to listen for exit signals
	go func() {
		sig := <-sigChan
		fmt.Println("\nReceived signal:", sig)
		cancel() // Cancel the context
	}()

	fmt.Println("App started. Press Ctrl+C to exit.")

	progressCh := make(chan bool, 5)
	defer close(progressCh)

	// Simulate work until context is canceled
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping work...")
				return
			default:
				progressCh <- true
				fmt.Println("start Working...")
				time.Sleep(5 * time.Second)
				fmt.Println("done Working...")
				<-progressCh
			}
		}
	}()

	<-ctx.Done()
	for i := 0; i < cap(progressCh); i++ {
		progressCh <- true
	}

	fmt.Println("Shutdown complete.")
}
