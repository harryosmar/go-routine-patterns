package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	fmt.Println("App started. Press Ctrl+C to exit.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("completed without exit signal")
	case <-sigChan:
		fmt.Println("exit signal received")
	}
}
