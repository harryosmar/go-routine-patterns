package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to handle exit signal
	go func() {
		sig := <-sigChan
		fmt.Println("\nReceived signal:", sig)
		cancel()
	}()

	fmt.Println("App started. Press Ctrl+C to exit.")

	var wg sync.WaitGroup
	wg.Add(2)

	// Worker 1
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping work1...")
				return
			default:
				fmt.Println("start Working1...")
				time.Sleep(5 * time.Second)
				fmt.Println("done Working1...")
			}
		}
	}()

	// Worker 2
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping work2...")
				return
			default:
				fmt.Println("start Working2...")
				time.Sleep(7 * time.Second)
				fmt.Println("done Working2...")
			}
		}
	}()

	// Wait for cancellation
	<-ctx.Done()
	fmt.Println("Waiting for tasks to finish...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All tasks finished.")
	case <-shutdownCtx.Done():
		fmt.Println("Shutdown timed out!")
	}

	fmt.Println("Shutdown complete.")
}
