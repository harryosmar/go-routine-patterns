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

	// Simulate work until context is canceled
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping work...")
				return
			default:
				fmt.Println("Working...")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Wait for cancellation
	<-ctx.Done()

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	fmt.Println("Shutting down gracefully...")
	cleanup(shutdownCtx)
	fmt.Println("Shutdown complete.")
}

// cleanup simulates resource cleanup (e.g., DB close, file flush)
func cleanup(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		// Simulate cleanup work
		time.Sleep(3 * time.Second)
		fmt.Println("Resources cleaned up.")
		close(done)
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Cleanup timed out!")
	case <-done:
		fmt.Println("Cleanup complete!")
	}
}
