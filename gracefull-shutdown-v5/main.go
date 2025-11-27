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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Println("\nReceived signal:", sig)
		cancel()
	}()

	fmt.Println("App started. Press Ctrl+C to exit.")

	tasks := make(chan int, 20)
	for i := 1; i <= 20; i++ {
		tasks <- i
	}
	close(tasks)

	var wg sync.WaitGroup
	workerCount := 5

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d stopping...\n", id)
					return
				case task, ok := <-tasks:
					if !ok {
						fmt.Printf("Worker %d done receiving all tasks, maybe some still processing...\n", id)
						return
					}
					fmt.Printf("Worker %d start processing task %d...\n", id, task)
					time.Sleep(5 * time.Second) // simulate work
					fmt.Printf("Worker %d done processing task %d...\n", id, task)
				}
			}
		}(w)
	}

	<-ctx.Done()
	fmt.Println("Waiting for workers to finish...")

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
		fmt.Println("All workers finished.")
	case <-shutdownCtx.Done():
		fmt.Println("Shutdown timed out!")
	}

	fmt.Println("Shutdown complete.")
}
