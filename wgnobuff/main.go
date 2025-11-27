package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	output1 := make(chan int)
	output2 := make(chan string)

	// read
	var s string
	var i int
	var readWg sync.WaitGroup
	readWg.Add(2)
	go func() {
		defer readWg.Done()
		s = <-output2
	}()
	go func() {
		defer readWg.Done()
		i = <-output1
	}()

	// write
	var writeWg sync.WaitGroup
	writeWg.Add(2)
	go func() {
		defer writeWg.Done()
		time.Sleep(1 * time.Second)
		output1 <- 1
		close(output1)
	}()

	go func() {
		defer writeWg.Done()
		time.Sleep(3 * time.Second)
		output2 <- "two"
		close(output2)
	}()

	writeWg.Wait()
	readWg.Wait()
	fmt.Println(i, s)
}
