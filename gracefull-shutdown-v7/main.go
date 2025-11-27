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

type Task struct {
	ID      int
	Retries int
}

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

	tasks := make(chan Task, 50)
	errors := make(chan error, 10)
	var wg sync.WaitGroup
	workerCount := 5
	maxRetries := 3

	// Start workers
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
						fmt.Printf("Worker %d finished all tasks.\n", id)
						return
					}
					fmt.Printf("Worker %d processing task %d...\n", id, task.ID)
					err := processTask(task)
					if err != nil {
						fmt.Printf("Worker %d error on task %d: %v\n", id, task.ID, err)
						if task.Retries < maxRetries {
							task.Retries++
							// Safe retry: check ctx before sending
							select {
							case <-ctx.Done():
								fmt.Printf("Shutdown started, skipping retry for task %d\n", task.ID)
							case tasks <- task:
								fmt.Printf("Retrying task %d (attempt %d)\n", task.ID, task.Retries)
							}
						} else {
							errors <- fmt.Errorf("task %d failed after %d retries", task.ID, maxRetries)
						}
					} else {
						fmt.Printf("Worker %d completed task %d.\n", id, task.ID)
					}
				}
			}
		}(w)
	}

	// Dynamic task producer
	go func() {
		for i := 1; i <= 20; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Producer stopping...")
				close(tasks)
				return
			case tasks <- Task{ID: i}:
				fmt.Printf("Produced task %d\n", i)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Error logger
	go func() {
		for err := range errors {
			fmt.Println("Error:", err)
		}
	}()

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

	close(errors)
	fmt.Println("Shutdown complete.")
}

// Simulate task processing with random failure
func processTask(task Task) error {
	time.Sleep(2 * time.Second)
	if task.ID%7 == 0 { // simulate error for some tasks
		return fmt.Errorf("simulated error on task %d", task.ID)
	}
	return nil
}
