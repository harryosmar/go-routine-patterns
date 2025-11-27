package main

import (
	"context"
	"fmt"
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

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("completed")
	case <-ctx.Done():
		fmt.Println("context cancelled")
	}
}
