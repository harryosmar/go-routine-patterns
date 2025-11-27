package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	output1 := make(chan int, 1)
	output2 := make(chan string, 1)

	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		output1 <- 1
		close(output1)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		output2 <- "two"
		close(output2)
	}()

	wg.Wait()
	fmt.Println(<-output1, <-output2)
}
