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

	var wg sync.WaitGroup
	wg.Add(2)

	// Simulate work 1 until context is canceled
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

	// Simulate work 2 until context is canceled
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

	<-ctx.Done()
	fmt.Println("Waiting for task to finish...")
	wg.Wait()

	fmt.Println("Shutdown complete.")
}
